package avator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	//avator "text-to-picture/services/auth_s/avator"         // 引入你的Avatar相关模块
	middlewire "text-to-picture/middlewire/jwt"
	"text-to-picture/models/user"

	db "text-to-picture/models/init"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

//	go test ./avator  或 go test -v ./avator（含测试详情）

// MockUserRepository 是一个模拟的用户仓库
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
	yamlFile, err := os.ReadFile("D:/go project/src/gocode/software-engineering/backend/text-to-picture/config/DBconfig/DBconfig.yaml")
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

// SetupRouter 设置 Gin 路由
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/auth/getavator",middlewire.JWTAuthMiddleware(), GetAvator)
	return r
}

// TestGetAvator_Success 测试成功获取头像的情况
func TestGetAvator_Success(t *testing.T) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 创建一个新的路由器
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

	// 创建一个 GET 请求
	req, _ := http.NewRequest("GET", "/auth/getavator", nil)
	req.Header.Set("Authorization", tokenString)

	// 创建一个响应记录器
	w := httptest.NewRecorder()

	// 执行请求
	router.ServeHTTP(w, req)

	fmt.Printf("\nTestGetAvator_Success-------------------Response Body: %s\n", w.Body.String())

	// 检查响应状态码
	assert.Equal(t, http.StatusOK, w.Code)

	// 检查响应体
	expectedResponse := AvatorResponse{
		Code: Success,
		Msg:  "获取头像成功",
		Data: "https://example.com/new-avatar.jpg", 
	}
	var actualResponse AvatorResponse
	json.Unmarshal(w.Body.Bytes(), &actualResponse)
	assert.Equal(t, expectedResponse, actualResponse)
}

// TestGetAvator_NoToken 测试请求头中缺少Token的情况
func TestGetAvator_NoToken(t *testing.T) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 创建一个新的路由器
	router := SetupRouter()

	// 创建一个 GET 请求，故意不设置Authorization头
	req, _ := http.NewRequest("GET", "/auth/getavator", nil)

	// 创建一个响应记录器
	w := httptest.NewRecorder()

	// 执行请求
	router.ServeHTTP(w, req)

	fmt.Printf("\nTestGetAvator_NoToken-------------------Response Body: %s\n", w.Body.String())

	// 检查响应状态码
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// 检查响应体
	expectedResponse := gin.H{
		"code":    float64(401),
		"message": "请求头中缺少Token",
	}
	var actualResponse gin.H
	json.Unmarshal(w.Body.Bytes(), &actualResponse)
	assert.Equal(t, expectedResponse, actualResponse)
}

// TestGetAvator_InvalidToken 测试无效Token的情况
func TestGetAvator_InvalidToken(t *testing.T) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 创建一个新的路由器
	router := SetupRouter()

	// 创建一个无效的Token
	tokenString := "invalid-token"

	// 创建一个 GET 请求
	req, _ := http.NewRequest("GET", "/auth/getavator", nil)
	req.Header.Set("Authorization", tokenString)

	// 创建一个响应记录器
	w := httptest.NewRecorder()

	// 执行请求
	router.ServeHTTP(w, req)

	fmt.Printf("\nTestGetAvator_InvalidToken-------------------Response Body: %s\n", w.Body.String())

	// 检查响应状态码
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// 检查响应体
	expectedResponse := gin.H{
		"code":    float64(401),
		"message": "无效的Token",
	}
	var actualResponse gin.H
	json.Unmarshal(w.Body.Bytes(), &actualResponse)
	assert.Equal(t, expectedResponse, actualResponse)
}

// TestGetAvator_ExpiredToken 测试Token过期的情况
func TestGetAvator_ExpiredToken(t *testing.T) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 创建一个新的路由器
	router := SetupRouter()

	// 创建一个过期的Token
	claims := &middlewire.Claims{
		Username: "czh1",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(-time.Hour).Unix(), // 设置过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(middlewire.JwtKey)

	// 创建一个 GET 请求
	req, _ := http.NewRequest("GET", "/auth/getavator", nil)
	req.Header.Set("Authorization", tokenString)

	// 创建一个响应记录器
	w := httptest.NewRecorder()

	// 执行请求
	router.ServeHTTP(w, req)

	fmt.Printf("\nTestGetAvator_ExpiredToken-------------------Response Body: %s\n", w.Body.String())

	// 检查响应状态码
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// 检查响应体
	expectedResponse := AvatorResponse{
		Code: Unauthorized,
		Msg:  "Token已过期",
		Data: nil,
	}
	var actualResponse AvatorResponse
	json.Unmarshal(w.Body.Bytes(), &actualResponse)
	assert.Equal(t, expectedResponse, actualResponse)
}

// TestGetAvator_ExpiredToken 测试名字解析出错的情况
func TestGetAvator_NameParseError(t *testing.T) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 创建一个新的路由器
	router := SetupRouter()

	// 创建一个有效的Token，但Username为空
	claims := &middlewire.Claims{
		Username: "",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1*time.Hour).Unix(), 
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(middlewire.JwtKey)

	// 创建一个 GET 请求
	req, _ := http.NewRequest("GET", "/auth/getavator", nil)
	req.Header.Set("Authorization", tokenString)

	// 创建一个响应记录器
	w := httptest.NewRecorder()

	// 执行请求
	router.ServeHTTP(w, req)

	fmt.Printf("\nTestGetAvator_NameParseError-------------------Response Body: %s\n", w.Body.String())

	// 检查响应状态码
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// 检查响应体
	expectedResponse := AvatorResponse{
		Code: Unauthorized,
		Msg:  "名字解析出错",
		Data: nil,
	}
	var actualResponse AvatorResponse
	json.Unmarshal(w.Body.Bytes(), &actualResponse)
	assert.Equal(t, expectedResponse, actualResponse)
}
