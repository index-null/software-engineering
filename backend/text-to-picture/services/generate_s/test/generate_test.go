package test_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"text-to-picture/api/generate"
	getDB "text-to-picture/config"
	middlewire "text-to-picture/middlewire/jwt"
	"text-to-picture/models/init"
	db "text-to-picture/models/init"
	"text-to-picture/models/repository/user_r"
	"text-to-picture/models/user"
	userLogin "text-to-picture/models/user"
	image "text-to-picture/services/generate_s"
	"time"
)

type MockUserRepository struct{}

// GetUserByName 是一个模拟的获取用户方法
func (m *MockUserRepository) GetUserByName(name string) (*user.UserInformation, error) {
	// 返回一个模拟的用户信息
	return &user.UserInformation{
		UserName:   name,
		Avatar_url: "http://example.com/avatar.png",
	}, nil
}

// DBConfig 结构体定义
type DBConfig struct {
	DB struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"db"`
}

// TestMain 是测试的入口函数
func TestMain(m *testing.M) {
	// 读取测试数据库配置
	yamlFile, err := os.ReadFile(getDB.GetDBConfigPath())
	if err != nil {
		fmt.Printf("Error reading DBconfig_test.yaml file: %v\n", err)
		os.Exit(1)
	}

	// 解析配置文件
	var dbconfig DBConfig
	err = yaml.Unmarshal(yamlFile, &dbconfig)
	if err != nil {
		fmt.Printf("Error parsing DBconfig_test.yaml file: %v\n", err)
		os.Exit(1)
	}

	// 设置数据库连接的环境变量
	os.Setenv("DB_USER", dbconfig.DB.User)
	os.Setenv("DB_PASSWORD", dbconfig.DB.Password)
	os.Setenv("DB_HOST", dbconfig.DB.Host)
	os.Setenv("DB_PORT", dbconfig.DB.Port)
	os.Setenv("DB_NAME", dbconfig.DB.Name)

	// 连接数据库
	if err := db.ConnectDatabase(); err != nil {
		fmt.Printf("Failed to connect to test database: %v\n", err)
		os.Exit(1)
	}

	// 运行测试
	code := m.Run()

	// 清理工作
	// 如有需要，在这里关闭数据库连接等

	os.Exit(code)
}

// SetupRouter sets up the Gin router for testing
func SetupRouter() *gin.Engine {
	imgGen := generate.NewImageGenerator()
	r := gin.Default()
	r.POST("/generate", func(c *gin.Context) {
		imgGen.ReturnImage(c)
	})
	return r
}

// TestGenerateImage_Success tests the successful case of generating an image
func TestGenerateImage_Success(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	// 创建一个有效的Token
	claims := &middlewire.Claims{
		Username: "testuser", //根据自己数据库已有的用户
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(), // 设置有效的过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(middlewire.JwtKey)
	user := userLogin.UserInformation{
		UserName: "testuser",
		Token:    tokenString,
		Email:    "test@qq.com",
		Password: "aaaaaa",
	}
	err := user_r.InsertUserInformation(models.DB, &user)
	if err != nil {
		log.Printf("插入错误%v", err)
	}
	// Set up the router
	router := SetupRouter()

	// Create valid input parameters
	input := image.ImageParaments{
		Prompt: "A beautiful sunset",
		Width:  1024,
		Height: 1024,
		Steps:  50,
		Seed:   123456,
	}

	// Marshal input into JSON
	jsonData, _ := json.Marshal(input)

	// Create a POST request
	req, _ := http.NewRequest("POST", "/generate", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenString)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Check response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Check response body
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "Image generated successfully", response["msg"])
}
func TestGenerateImage_NoToken(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Set up the router
	router := SetupRouter()

	// Create invalid input parameters (width out of range)
	input := image.ImageParaments{
		Prompt: "A beautiful sunset",
		Width:  1024,
		Height: 1024,
		Steps:  50,
		Seed:   123456,
	}

	// Marshal input into JSON
	jsonData, _ := json.Marshal(input)

	// Create a POST request
	req, _ := http.NewRequest("POST", "/generate", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Check response status code
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Check response body
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(400), response["code"])
	assert.Contains(t, response["msg"], "宽度不在范围内")
}

// TestGenerateImage_InvalidParameters tests the case where invalid parameters are provided
func TestGenerateImage_InvalidParameters(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Set up the router
	router := SetupRouter()
	// 创建一个有效的Token
	claims := &middlewire.Claims{
		Username: "testuser", //根据自己数据库已有的用户
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(), // 设置有效的过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(middlewire.JwtKey)
	// Create invalid input parameters (width out of range)
	input := map[string]interface{}{
		"prompt": "A beautiful sunset",
		"width":  2000,
		"height": 512,
		"steps":  50,
		"seed":   123456,
	}

	// Marshal input into JSON
	jsonData, _ := json.Marshal(input)

	// Create a POST request
	req, _ := http.NewRequest("POST", "/generate", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenString)
	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Check response status code
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Check response body
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(400), response["code"])
	assert.Contains(t, response["msg"], "宽度不在范围内")
}

// TestGenerateImage_MissingParameters tests the case where required parameters are missing
func TestGenerateImage_MissingParameters(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Set up the router
	router := SetupRouter()
	// 创建一个有效的Token
	claims := &middlewire.Claims{
		Username: "testuser", //根据自己数据库已有的用户
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(), // 设置有效的过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(middlewire.JwtKey)
	// Create input with missing "prompt"
	input := map[string]interface{}{
		"width":  512,
		"height": 512,
		"steps":  50,
		"seed":   123456,
	}

	// Marshal input into JSON
	jsonData, _ := json.Marshal(input)

	// Create a POST request
	req, _ := http.NewRequest("POST", "/generate", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenString)
	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Check response status code
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Check response body
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(400), response["code"])
	assert.Contains(t, response["msg"], "缺乏提示词")
}

// TestGenerateImage_ServerError tests the case where the server encounters an error
func TestGenerateImage_ServerError(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	// 创建一个有效的Token
	claims := &middlewire.Claims{
		Username: "testuser", //根据自己数据库已有的用户
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(), // 设置有效的过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(middlewire.JwtKey)
	// Simulate database connection failure
	// (you can mock database connection here to force an error)

	// Set up the router
	router := SetupRouter()

	// Create valid input parameters
	input := image.ImageParaments{
		Prompt: "A beautiful sunset",
		Width:  512,
		Height: 512,
		Steps:  50,
		Seed:   123456,
	}

	// Marshal input into JSON
	jsonData, _ := json.Marshal(input)

	// Create a POST request
	req, _ := http.NewRequest("POST", "/generate", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenString)
	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Check response status code
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Check response body
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(500), response["code"])
	assert.Equal(t, "Server error", response["msg"])
}
