package generate_s

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	db "text-to-picture/models/init"
	u "text-to-picture/models/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockImageParaments 是一个模拟的文生图请求参数结构体
type MockImageParaments struct {
	Prompt         string `json:"prompt"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	Steps          int    `json:"steps"`
	SamplingMethod string `json:"sampling_method"`
	Seed           string `json:"seed"`
}

// MockImageResponse 是一个模拟的文生图响应结构体
type MockImageResponse struct {
	Code      int    `json:"code"`
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Image_url string `json:"image_url"`
}

func TestGenerateImageValidRequest(t *testing.T) {
	// 初始化 Gin 路由
	r := gin.Default()

	generator := &ImageGeneratorImpl{}

	r.POST("/generate", generator.ReturnImage)

	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "07080031")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "postgres")
	db.ConnectDatabase()
	db.InitDB()

	// 创建一个测试服务器
	ts := httptest.NewServer(r)
	defer ts.Close()

	// 模拟请求数据
	mockRequest := &MockImageParaments{
		Prompt:         "生成一张美丽的风景画",
		Width:          512,
		Height:         512,
		Steps:          50,
		SamplingMethod: "DDIM",
		Seed:           "12345",
	}

	// 将请求数据转换为 JSON
	jsonData, _ := json.Marshal(mockRequest)

	// 发送 POST 请求
	req, _ := http.NewRequest("POST", ts.URL+"/generate", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)

	// 读取响应体
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	// 解析响应体
	var response MockImageResponse
	json.Unmarshal(body, &response)

	// 断言响应状态
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, true, response.Success)
	assert.NotEmpty(t, response.Image_url)
}
func TestGenerateImageSuccessfulRequest(t *testing.T) {
	// 初始化 Gin 路由
	r := gin.Default()

	generator := &ImageGeneratorImpl{}
	r.POST("/generate", generator.ReturnImage)

	db.ConnectDatabase()
	// 创建一个测试服务器
	ts := httptest.NewServer(r)
	defer ts.Close()

	// 模拟请求数据
	mockRequest := &MockImageParaments{
		Prompt:         "生成一张美丽的风景画",
		Width:          512,
		Height:         512,
		Steps:          50,
		SamplingMethod: "DDIM",
		Seed:           "12345",
	}

	// 将请求数据转换为 JSON
	jsonData, _ := json.Marshal(mockRequest)

	// 发送 POST 请求
	req, _ := http.NewRequest("POST", ts.URL+"/generate", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)

	// 读取响应体
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	// 解析响应体
	var response MockImageResponse
	json.Unmarshal(body, &response)

	// 断言响应状态
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, true, response.Success)
	assert.NotEmpty(t, response.Image_url)
}

func TestGenerateImageLackingKeyWordRequest(t *testing.T) {
	// 初始化 Gin 路由
	r := gin.Default()
	generator := &ImageGeneratorImpl{}
	r.POST("/generate", generator.ReturnImage)

	db.ConnectDatabase()

	// 创建一个测试服务器
	ts := httptest.NewServer(r)
	defer ts.Close()

	// 模拟无效的请求数据
	mockRequest := &MockImageParaments{
		Prompt:         "", // 缺少提示词
		Width:          512,
		Height:         512,
		Steps:          50,
		SamplingMethod: "DDIM",
		Seed:           "12345",
	}

	// 将请求数据转换为 JSON
	jsonData, _ := json.Marshal(mockRequest)

	// 发送 POST 请求
	req, _ := http.NewRequest("POST", ts.URL+"/generate", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)

	// 读取响应体
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	// 解析响应体
	var response MockImageResponse
	json.Unmarshal(body, &response)

	// 断言响应状态
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, "缺乏提示词", response.Message)
}

func TestGenerateImageInternalError(t *testing.T) {
	// 初始化 Gin 路由
	r := gin.Default()

	generator := &ImageGeneratorImpl{}

	r.POST("/generate", func(c *gin.Context) {
		generator.ReturnImage(c)
	})

	db.ConnectDatabase()

	// 创建一个测试服务器
	ts := httptest.NewServer(r)
	defer ts.Close()

	// 模拟请求数据
	mockRequest := &MockImageParaments{
		Prompt:         "生成一张美丽的风景画",
		Width:          512,
		Height:         512,
		Steps:          50,
		SamplingMethod: "DDIM",
		Seed:           "12345",
	}

	// 将请求数据转换为 JSON
	jsonData, _ := json.Marshal(mockRequest)

	// 发送 POST 请求
	req, _ := http.NewRequest("POST", ts.URL+"/generate", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)

	// 读取响应体
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	// 解析响应体
	var response MockImageResponse
	json.Unmarshal(body, &response)

	// 断言响应状态
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, 500, response.Code)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, "图片生成失败", response.Message)
}

func TestGenerateImageWidthTooSmallError(t *testing.T) {
	// 初始化 Gin 路由
	r := gin.Default()

	generator := &ImageGeneratorImpl{}

	r.POST("/generate", generator.ReturnImage)

	db.ConnectDatabase()

	// 创建一个测试服务器
	ts := httptest.NewServer(r)
	defer ts.Close()

	// 模拟请求数据
	mockRequest := &MockImageParaments{
		Prompt:         "生成一张美丽的风景画",
		Width:          64,
		Height:         512,
		Steps:          50,
		SamplingMethod: "DDIM",
		Seed:           "12345",
	}

	// 将请求数据转换为 JSON
	jsonData, _ := json.Marshal(mockRequest)

	// 发送 POST 请求
	req, _ := http.NewRequest("POST", ts.URL+"/generate", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)

	// 读取响应体
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	// 解析响应体
	var response MockImageResponse
	json.Unmarshal(body, &response)

	// 断言响应状态
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, false, response.Success)
	assert.NotEmpty(t, "宽度不在范围内", response.Message)
}

func TestGenerateImageHeightTooBigError(t *testing.T) {
	// 初始化 Gin 路由
	r := gin.Default()

	generator := &ImageGeneratorImpl{}

	r.POST("/generate", generator.ReturnImage)

	db.ConnectDatabase()

	// 创建一个测试服务器
	ts := httptest.NewServer(r)
	defer ts.Close()

	// 模拟请求数据
	mockRequest := &MockImageParaments{
		Prompt:         "生成一张美丽的风景画",
		Width:          512,
		Height:         2048,
		Steps:          50,
		SamplingMethod: "DDIM",
		Seed:           "12345",
	}

	// 将请求数据转换为 JSON
	jsonData, _ := json.Marshal(mockRequest)

	// 发送 POST 请求
	req, _ := http.NewRequest("POST", ts.URL+"/generate", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)

	// 读取响应体
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	// 解析响应体
	var response MockImageResponse
	json.Unmarshal(body, &response)

	// 断言响应状态
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, false, response.Success)
	assert.NotEmpty(t, "高度不在范围内", response.Message)
}

func TestGenerateImageStepTooBigError(t *testing.T) {
	// 初始化 Gin 路由
	r := gin.Default()

	generator := &ImageGeneratorImpl{}

	r.POST("/generate", generator.ReturnImage)

	db.ConnectDatabase()

	// 创建一个测试服务器
	ts := httptest.NewServer(r)
	defer ts.Close()

	// 模拟请求数据
	mockRequest := &MockImageParaments{
		Prompt:         "生成一张美丽的风景画",
		Width:          512,
		Height:         512,
		Steps:          130,
		SamplingMethod: "DDIM",
		Seed:           "12345",
	}

	// 将请求数据转换为 JSON
	jsonData, _ := json.Marshal(mockRequest)

	// 发送 POST 请求
	req, _ := http.NewRequest("POST", ts.URL+"/generate", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)

	// 读取响应体
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	// 解析响应体
	var response MockImageResponse
	json.Unmarshal(body, &response)

	// 断言响应状态
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, false, response.Success)
	assert.NotEmpty(t, "步数不在范围内", response.Message)
}
func TestGenerateImageUserScoreDecrease(t *testing.T) {
	// 初始化 Gin 路由
	r := gin.Default()

	generator := &ImageGeneratorImpl{}
	r.POST("/generate", generator.ReturnImage)

	db.ConnectDatabase()

	// 创建一个测试服务器
	ts := httptest.NewServer(r)
	defer ts.Close()

	// 模拟请求数据
	mockRequest := &MockImageParaments{
		Prompt:         "生成一张美丽的风景画",
		Width:          512,
		Height:         512,
		Steps:          50,
		SamplingMethod: "DDIM",
		Seed:           "12345",
	}

	// 将请求数据转换为 JSON
	jsonData, _ := json.Marshal(mockRequest)

	// 准备用户数据
	user := u.UserInformation{
		UserName: "testuser",
		Score:    50, // 初始积分大于20
	}
	db.DB.Create(&user)

	// 发送 POST 请求
	req, _ := http.NewRequest("POST", ts.URL+"/generate", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer testtoken") // 假设存在这样的认证机制

	client := &http.Client{}
	resp, _ := client.Do(req)

	// 读取响应体
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	// 解析响应体
	var response MockImageResponse
	json.Unmarshal(body, &response)

	// 断言响应状态
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, true, response.Success)
	assert.NotEmpty(t, response.Image_url)

	// 验证用户积分是否减少
	var updatedUser u.UserInformation
	db.DB.Where("username = ?", "testuser").First(&updatedUser)
	assert.Equal(t, 30, updatedUser.Score) // 初始50减去20

	// 清理测试数据
	db.DB.Unscoped().Delete(&user)
}
func TestGenerateImageScoreRecordCreated(t *testing.T) {
	// 初始化 Gin 路由
	r := gin.Default()

	generator := &ImageGeneratorImpl{}
	r.POST("/generate", generator.ReturnImage)

	db.ConnectDatabase()

	// 创建一个测试服务器
	ts := httptest.NewServer(r)
	defer ts.Close()

	// 模拟请求数据
	mockRequest := &MockImageParaments{
		Prompt:         "生成一张美丽的风景画",
		Width:          512,
		Height:         512,
		Steps:          50,
		SamplingMethod: "DDIM",
		Seed:           "12345",
	}

	// 将请求数据转换为 JSON
	jsonData, _ := json.Marshal(mockRequest)

	// 准备用户数据
	user := u.UserInformation{
		UserName: "testuser",
		Score:    50, // 初始积分大于20
	}
	db.DB.Create(&user)

	// 发送 POST 请求
	req, _ := http.NewRequest("POST", ts.URL+"/generate", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer testtoken") // 假设存在这样的认证机制

	client := &http.Client{}
	resp, _ := client.Do(req)

	// 读取响应体
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	// 解析响应体
	var response MockImageResponse
	json.Unmarshal(body, &response)

	// 断言响应状态
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, true, response.Success)
	assert.NotEmpty(t, response.Image_url)

	// 验证积分记录是否创建
	var record u.UserScore
	db.DB.Where("username = ? AND record = ?", "testuser", "积分-20").First(&record)
	assert.NotEmpty(t, record.ID) // 记录存在

	// 清理测试数据
	db.DB.Unscoped().Delete(&user)
	db.DB.Unscoped().Delete(&record)
}

func TestGenerateImageInsufficientScore(t *testing.T) {
	// 初始化 Gin 路由
	r := gin.Default()

	generator := &ImageGeneratorImpl{}
	r.POST("/generate", generator.ReturnImage)

	db.ConnectDatabase()

	// 创建一个测试服务器
	ts := httptest.NewServer(r)
	defer ts.Close()

	// 模拟请求数据
	mockRequest := &MockImageParaments{
		Prompt:         "生成一张美丽的风景画",
		Width:          512,
		Height:         512,
		Steps:          50,
		SamplingMethod: "DDIM",
		Seed:           "12345",
	}

	// 将请求数据转换为 JSON
	jsonData, _ := json.Marshal(mockRequest)

	// 准备用户数据
	user := u.UserInformation{
		UserName: "testuser",
		Score:    10, // 初始积分小于20
	}
	db.DB.Create(&user)

	// 发送 POST 请求
	req, _ := http.NewRequest("POST", ts.URL+"/generate", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer testtoken") // 假设存在这样的认证机制

	client := &http.Client{}
	resp, _ := client.Do(req)

	// 读取响应体
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	// 解析响应体
	var response MockImageResponse
	json.Unmarshal(body, &response)

	// 断言响应状态
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.Equal(t, 401, response.Code)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, "用户积分不足", response.Message)

	// 验证用户积分未改变
	var unchangedUser u.UserInformation
	db.DB.Where("username = ?", "testuser").First(&unchangedUser)
	assert.Equal(t, 10, unchangedUser.Score) // 积分未变化

	// 清理测试数据
	db.DB.Unscoped().Delete(&user)
}

func TestGenerateImageNoSuchMethodError(t *testing.T) {
	// 初始化 Gin 路由
	r := gin.Default()

	generator := &ImageGeneratorImpl{}

	r.POST("/generate", generator.ReturnImage)

	db.ConnectDatabase()

	// 创建一个测试服务器
	ts := httptest.NewServer(r)
	defer ts.Close()

	// 模拟请求数据
	mockRequest := &MockImageParaments{
		Prompt:         "生成一张美丽的风景画",
		Width:          512,
		Height:         512,
		Steps:          50,
		SamplingMethod: "AABC",
		Seed:           "12345",
	}

	// 将请求数据转换为 JSON
	jsonData, _ := json.Marshal(mockRequest)

	// 发送 POST 请求
	req, _ := http.NewRequest("POST", ts.URL+"/generate", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)

	// 读取响应体
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	// 解析响应体
	var response MockImageResponse
	json.Unmarshal(body, &response)

	// 断言响应状态
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, false, response.Success)
	assert.NotEmpty(t, "采样方法不在范围内", response.Message)
}

func TestGenerateImageLackingSeedError(t *testing.T) {
	// 初始化 Gin 路由
	r := gin.Default()

	generator := &ImageGeneratorImpl{}

	r.POST("/generate", generator.ReturnImage)

	db.ConnectDatabase()

	// 创建一个测试服务器
	ts := httptest.NewServer(r)
	defer ts.Close()

	// 模拟请求数据
	mockRequest := &MockImageParaments{
		Prompt:         "生成一张美丽的风景画",
		Width:          512,
		Height:         512,
		Steps:          50,
		SamplingMethod: "DDIM",
		Seed:           "",
	}

	// 将请求数据转换为 JSON
	jsonData, _ := json.Marshal(mockRequest)

	// 发送 POST 请求
	req, _ := http.NewRequest("POST", ts.URL+"/generate", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)

	// 读取响应体
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	// 解析响应体
	var response MockImageResponse
	json.Unmarshal(body, &response)

	// 断言响应状态
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, false, response.Success)
	assert.NotEmpty(t, "缺乏种子", response.Message)
}
