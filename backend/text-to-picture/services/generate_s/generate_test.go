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

	os.Setenv("DB_USER", "YourUsername")
	os.Setenv("DB_PASSWORD", "YourPassword")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "YourDBName")
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
