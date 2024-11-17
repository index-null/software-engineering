package generate_s

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"log"
	"net/url"
	"os"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
)

// 传入的图片参数
// @name ImageParaments
type ImageParaments struct {
	Prompt         string `json:"prompt" binding:"required" fault:"缺乏提示词"`
	Width          int    `json:"width" binding:"required,min=128,max=1024" fault:"宽度不在范围内"`
	Height         int    `json:"height" binding:"required,min=128,max=1024" fault:"高度不在范围内"`
	Steps          int    `json:"steps" binding:"required,min=1,max=100" fault:"步数不在范围内"`
	SamplingMethod string `json:"sampling_method" binding:"required,oneof=DDIM PLMS K-LMS" fault:"采样方法不在范围内"`
	Seed           string `json:"seed" binding:"required," fault:"缺乏种子"`
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
func IsValidURL(inputURL string) bool {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return false
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	}

	return true
}

var imageParaments ImageParaments

// 接收传来的图片参数,并进行数据校验
func AcceptParaments(c *gin.Context) error {
	if err := c.ShouldBindJSON(&imageParaments); err != nil {
		fault := ParamentsError(err, &imageParaments)
		return fmt.Errorf(fault)
	}

	return nil
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
// @Router /generate [post]
func (*ImageGeneratorImpl) ReturnImage(c *gin.Context) {
	//校验参数
	if err := AcceptParaments(c); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": err,
		})
		return
	}

	imageUrl, err := GenerateImage()

	//校验生成图片
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "图片生成失败",
		})
		return
	}

	//校验图片url
	if !IsValidURL(imageUrl) {
		c.JSON(500, gin.H{
			"success": false,
			"message": "无效url",
		})
		return
	}

	c.JSON(200, gin.H{
		"success":   true,
		"image_url": imageUrl,
	})
}

func GenerateImage() (string, error) {
	//这里把图片上传到OSS,OSS会那里返回包含图片URL的json
	urloss, err := SavetoOss()
	return urloss, err
}

var client *oss.Client // 全局变量用来存储OSS客户端实例
func SavetoOss() (string, error) {
	if err := godotenv.Load("config\\oss\\oss.env"); err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}
	// 从环境变量中获取访问凭证
	accessKeyID := os.Getenv("accessKeyId")
	accessKeySecret := os.Getenv("accessKeySecret")
	region := os.Getenv("region")
	bucketName := os.Getenv("bucket")
	// 从环境变量中获取访问凭证。运行本代码示例之前，请确保已设置环境变量OSS_ACCESS_KEY_ID和OSS_ACCESS_KEY_SECRET。
	provider, err := oss.NewEnvironmentVariableCredentialsProvider()
	if err != nil {
		log.Fatalf("Failed to create credentials provider: %v", err)
	}

	// 创建OSSClient实例。
	// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	// yourRegion填写Bucket所在地域，以华东1（杭州）为例，填写为cn-hangzhou。其它Region请按实际情况填写。
	clientOptions := []oss.ClientOption{oss.SetCredentialsProvider(&provider)}
	clientOptions = append(clientOptions, oss.Region(region))
	// 设置签名版本
	clientOptions = append(clientOptions, oss.AuthVersion(oss.AuthV4))
	client, err = oss.New("https://oss-"+region+".aliyuncs.com", accessKeyID, accessKeySecret, clientOptions...)
	if err != nil {
		log.Fatalf("Failed to create OSS client: %v", err)
	}
	// 填写存储空间名称，例如examplebucket。

	// 示例操作：上传文件。
	filetime := time.Now().Format("2024-01-01 00:00:00")
	encodedPrompt := url.QueryEscape(imageParaments.Prompt)
	objectName := "generate/" + encodedPrompt + "-" + filetime + ".png"
	localFileName := "E:/截图/新建文件夹/图片1.png" //测试就换成自己要上传的图片即可
	if err := uploadFile(bucketName, objectName, localFileName); err != nil {
		handleError(err)
	}
	return objectName, err
}

// handleError 用于处理不可恢复的错误，并记录错误信息后终止程序。
func handleError(err error) {
	log.Fatalf("Error: %v", err)
}

// uploadFile 用于将本地文件上传到OSS存储桶。
// 参数：
//
//	bucketName - 存储空间名称。
//	objectName - Object完整路径，完整路径中不包含Bucket名称。
//	localFileName - 本地文件的完整路径。
//	endpoint - Bucket对应的Endpoint。
//
// 如果成功，记录成功日志；否则，返回错误。
func uploadFile(bucketName, objectName, localFileName string) error {
	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	// 上传文件。
	err = bucket.PutObjectFromFile(objectName, localFileName)
	if err != nil {
		return err
	}

	// 文件上传成功后，记录日志。
	log.Printf("File uploaded successfully to %s/%s", bucketName, objectName)
	return nil
}
