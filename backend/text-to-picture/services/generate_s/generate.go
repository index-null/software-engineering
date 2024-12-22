package generate_s

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	//"net/http"
	// "net/url"
	"os"
	"reflect"
	i "text-to-picture/models/image"
	db "text-to-picture/models/init"
	u "text-to-picture/models/user"
	"time"

	"github.com/gin-gonic/gin"
	//文件路径操作包
	"path/filepath"
)

// 传入的图片参数
// @name ImageParaments
type ImageParaments struct {
	Prompt string `json:"prompt" binding:"required" fault:"缺乏提示词"`
	Width  int    `json:"width" binding:"required,min=128,max=1281" fault:"宽度不在范围内"`
	Height int    `json:"height" binding:"required,min=128,max=1281" fault:"高度不在范围内"`
	Steps  int    `json:"steps" binding:"required,min=1,max=100" fault:"步数不在范围内"`
	Seed   int    `json:"seed" binding:"required,min=1,max=4294967290" fault:"种子不在范围"`
}

// 获取在Tag中的fault信息
func ParamentsError(err error, obj any) string {
	getObj := reflect.TypeOf(obj)
	if err, ok := err.(validator.ValidationErrors); ok {
		for _, e := range err {
			if f, exists := getObj.Elem().FieldByName(e.Field()); exists {
				fault := f.Tag.Get("fault")
				return fault
			}
		}
	}
	return err.Error()
}

// 判断OSS返回的url是否正确
//func IsValidURL(inputURL string) bool {
//	parsedURL, err := url.Parse(inputURL)
//	if err != nil {
//		return false
//	}
//
//	if parsedURL.Scheme == "" || parsedURL.Host == "" {
//		return false
//	}
//
//	return true
//}

// 接收传来的图片参数,并进行数据校验
func AcceptParaments(c *gin.Context) (ImageParaments, error) {
	var imageParaments ImageParaments
	if err := c.ShouldBindJSON(&imageParaments); err != nil {
		fault := ParamentsError(err, &imageParaments)
		return imageParaments, fmt.Errorf(fault)
	}

	return imageParaments, nil
}

// ImageGeneratorImpl 实现了 ImageGenerator 接口
type ImageGeneratorImpl struct{}

// @Summary 生成图片
// @Description 根据传入的参数生成图片并返回图片的URL
// @Tags 图片生成
// @Accept json
// @Produce json
// @Param imageParaments body ImageParaments true "图片参数"
// @Success 200 {object} map[string]interface{} "成功响应"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 500 {object} map[string]interface{} "内部错误"
// @Router /auth/generate [post]
func (*ImageGeneratorImpl) ReturnImage(c *gin.Context) {
	var imageParaments ImageParaments
	var err error
	// 校验参数
	if imageParaments, err = AcceptParaments(c); err != nil {
		log.Printf("参数错误: %v", err)
		c.JSON(400, gin.H{
			"code":    400,
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 从上下文中获取用户名
	Username, exists := c.Get("username")
	if !exists {
		log.Printf("未找到用户名")
		c.JSON(401, gin.H{
			"code":    401,
			"success": false,
			"message": "未找到用户信息",
		})
		return
	}
	username := Username.(string)

	var user u.UserInformation
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		log.Printf("Failed to query user information: %v", err)
		c.JSON(401, gin.H{
			"code":    401,
			"success": false,
			"message": "用户信息查询失败",
		})
		return
	}

	if user.Score < 20 {
		c.JSON(401, gin.H{
			"code":    401,
			"success": false,
			"message": "用户积分不足",
		})
		return
	}

	user.Score -= 20
	if err := db.DB.Save(&user).Error; err != nil {
		log.Printf("Failed to update user score: %v", err)
		c.JSON(401, gin.H{
			"code":    401,
			"success": false,
			"message": "用户积分更新失败",
		})
		return
	}
	var record u.UserScore
	record.Username = username
	record.Record = "积分-20"
	if err := db.DB.Create(&record).Error; err != nil {
		log.Printf("Failed to create record: %v", err)
		c.JSON(401, gin.H{
			"code":    401,
			"success": false,
			"message": "积分记录创建失败",
		})
		return
	}
	// 生成图片并传递用户名
	imageUrl, err := GenerateImage(username, imageParaments)
	//校验生成图片
	if err != nil {
		log.Printf("图片生成失败: %v", err)
		c.JSON(500, gin.H{
			"code":    500,
			"success": false,
			"message": "图片生成失败",
		})
		return
	}

	////校验图片url
	//if !IsValidURL(imageUrl) {
	//	log.Printf("无效url: %v", err)
	//	c.JSON(500, gin.H{
	//		"success": false,
	//		"message": "无效url",
	//	})
	//	return
	//}
	msg := fmt.Sprintf("用户当前积分为%v", user.Score)
	c.JSON(200, gin.H{
		"code":      200,
		"success":   true,
		"image_url": imageUrl,
		"message":   msg,
	})
}

func GenerateImage(username string, imageParaments ImageParaments) (string, error) {
	//这里把图片上传到OSS,OSS会那里返回包含图片URL的json
	urloss, err := SavetoOss(imageParaments)
	if err != nil {
		return "", err
	}
	// 创建 ImageInformation 实例
	imageInfo := i.ImageInformation{
		UserName: username, // 实际使用时应该从会话信息中获取真实用户名
		Params: fmt.Sprintf("\"Prompt\": \"%s\", \"Width\": \"%d\", \"Height\": \"%d\", \"Steps\": \"%d\"",
			imageParaments.Prompt, imageParaments.Width, imageParaments.Height, imageParaments.Steps),
		Picture:     urloss, // 保存生成的图片 URL
		Create_time: time.Now(),
	}

	// 保存到数据库
	if err := db.DB.Create(&imageInfo).Error; err != nil {
		log.Printf("Failed to save image information to database: %v", err)
		return urloss, err
	}

	return urloss, err
}

var client *oss.Client // 全局变量用来存储OSS客户端实例
func SavetoOss(imageParaments ImageParaments) (string, error) {
	// 构建跨平台的路径
	localFileName := "assets/examples/images/3.jpg" //测试就换成自己要上传的图片即可

	envPath := filepath.Join("config", "oss", "oss.env")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Failed to load .env file: %v", err)
	}
	//localFileName, err := GenerateFromWebUI(imageParaments)
	// 从环境变量中获取访问凭证
	region := os.Getenv("OSS_REGION")
	bucketName := os.Getenv("OSS_BUCKET")
	accessKeyId := os.Getenv("OSS_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("OSS_ACCESS_KEY_SECRET")
	object := "generate"
	// 检查bucket名称是否为空
	if len(bucketName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, bucket name required")
	}

	// 检查region是否为空
	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

	// 检查object名称是否为空
	if len(object) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, object name required")
	}

	// 创建OSS客户端
	endPoint := region + ".aliyuncs.com"
	client, err := oss.New(endPoint, accessKeyId, accessKeySecret)
	bucket, err := client.Bucket(bucketName)
	file, err := os.Open(localFileName)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	objectName := "generate/" + localFileName // You can modify this path as needed
	result := bucket.PutObject(objectName, file)
	log.Printf("result: %v", result)
	url, err := bucket.SignURL(objectName, oss.HTTPGet, 3600)
	if err != nil {
		return "", fmt.Errorf("failed to generate URL: %v", err)
	}
	return url, err
}

// handleError 用于处理不可恢复的错误，并记录错误信息后终止程序。
func handleError(err error) {
	log.Printf("Error: %v", err)
}

func GenerateFromWebUI(imageParaments ImageParaments) (string, error) {
	var Url string
	apiKey := os.Getenv("GEN_API_KEY")
	type Parameters struct {
		Size string `json:"size"`
		Seed int    `json:"seed"`
	}
	type Input struct {
		Prompt string `json:"prompt"`
	}
	type RequestBody struct {
		Model      string     `json:"model"`
		Input      Input      `json:"input"`
		Parameters Parameters `json:"parameters"`
	}
	type TaskResponse struct {
		Output struct {
			TaskStatus string `json:"task_status"`
			TaskID     string `json:"task_id"`
		} `json:"output"`
		RequestID string `json:"request_id"`
	}
	type TaskStatusResponse struct {
		RequestID string `json:"request_id"`
		Output    struct {
			TaskStatus string `json:"task_status"`
			TaskID     string `json:"task_id"`
			Code       string `json:"code"`
			Message    string `json:"message"`
			Result     []struct {
				URL string `json:"url"`
			} `json:"results"`
		} `json:"output"`
	}
	// 创建 HTTP 客户端
	client := &http.Client{}
	size := fmt.Sprintf("%d*%d", imageParaments.Width, imageParaments.Height)
	// 构建请求体
	requestBody := RequestBody{
		Model: "wanx-v1",
		Input: Input{
			Prompt: imageParaments.Prompt,
		},
		Parameters: Parameters{
			//Style: "<auto>",
			Size: size,
			Seed: imageParaments.Seed,

			//N:    1,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// 创建 POST 请求来创建任务
	req, err := http.NewRequest("POST", "https://dashscope.aliyuncs.com/api/v1/services/aigc/text2image/image-synthesis", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-DashScope-Async", "enable")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应体
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析响应以获取任务ID
	var taskResponse TaskResponse
	err = json.Unmarshal(bodyText, &taskResponse)
	if err != nil {
		return "", fmt.Errorf("解析响应失败%v", err)
	}

	taskID := taskResponse.Output.TaskID
	if taskID == "" {
		return "", fmt.Errorf("任务ID为空，请检查请求是否成功%v", taskResponse)
	}
	// 轮询任务状态
	for {
		time.Sleep(30 * time.Second) // 每分钟轮询一次

		// 创建 GET 请求来查询任务状态
		statusReq, err := http.NewRequest("GET", fmt.Sprintf("https://dashscope.aliyuncs.com/api/v1/tasks/%s", taskID), nil)
		if err != nil {
			return "", err
		}

		// 设置请求头
		statusReq.Header.Set("Authorization", "Bearer "+apiKey)

		// 发送请求
		statusResp, err := client.Do(statusReq)
		if err != nil {
			return "", err
		}
		defer statusResp.Body.Close()

		// 读取响应体
		statusBodyText, err := io.ReadAll(statusResp.Body)
		if err != nil {
			return "", fmt.Errorf("读取响应失败%v", err)
		}

		// 解析响应以获取任务状态
		var taskStatusResponse TaskStatusResponse
		err = json.Unmarshal(statusBodyText, &taskStatusResponse)
		if err != nil {
			return "", fmt.Errorf("解析响应失败%v", err)
		}

		if taskStatusResponse.Output.TaskStatus == "SUCCEEDED" {
			log.Printf("生成成功，url为%v", taskStatusResponse.Output.Result[0].URL)
			url := taskStatusResponse.Output.Result[0].URL
			// 使用 filepath.Join 构建跨平台路径
			filetime := time.Now().Format("2006-01-02_15-04-05")
			localFileName := filepath.Join("assets/examples/images/", fmt.Sprintf("image_%s.png", filetime))
			if err := downloadImageFromWebUI(url, localFileName); err != nil {
				log.Printf("图片从模型保存到本地错误%v", err)
			}
			Url = localFileName
			break
		} else if taskStatusResponse.Output.TaskStatus == "FAILED" {
			log.Printf("任务失败%v,message:%v", taskStatusResponse.Output.Code, taskStatusResponse.Output.Message)
			return "", fmt.Errorf("任务失败%v,message:%v", taskStatusResponse.Output.Code, taskStatusResponse.Output.Message)
		}
	}
	return Url, nil
}

func downloadImageFromWebUI(url string, destinationPath string) error {
	// 发送HTTP GET请求
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("下载图片失败: %v", err)
	}
	defer resp.Body.Close()

	// 创建本地文件
	out, err := os.Create(destinationPath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer out.Close()

	// 复制数据
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("保存图片失败: %v", err)
	}

	return nil
}
