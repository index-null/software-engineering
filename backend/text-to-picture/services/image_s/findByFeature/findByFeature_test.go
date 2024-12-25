package findByFeature

import (
	"encoding/json"
	"fmt"
	//"log"
	"net/http"
	"net/http/httptest"
	"os"
	//"path/filepath"
	//"runtime"
	"testing"

	//middlewire "text-to-picture/middlewire/jwt"
	getDB "text-to-picture/config"
	middlewire "text-to-picture/middlewire/jwt"
	"text-to-picture/models/image"
	db "text-to-picture/models/init"
	"time"

	//"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

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
		fmt.Printf("Error reading configs.yaml file: %v\n", err)
		os.Exit(1)
	}

	// 解析配置文件
	var dbconfig DBConfig
	err = yaml.Unmarshal(yamlFile, &dbconfig)
	if err != nil {
		fmt.Printf("Error parsing configs.yaml file: %v\n", err)
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
	// 在这里关闭数据库连接等

	os.Exit(code)
}

// SetupRouter 设置 Gin 路由
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/auth/image/feature", middlewire.JWTAuthMiddleware(), FindByFeature)
	return r
}

// TestFindByFeature_Success 测试成功(isOwn不为true)获取图片的情况
func TestFindByFeature_Success(t *testing.T) {
	fmt.Println()
	fmt.Println()
	fmt.Println("\n----------------------------------TestFindByFeature_Success")
	gin.SetMode(gin.TestMode)
	router := SetupRouter()

	// 插入测试数据
	imageData := image.ImageInformation{
		UserName:    "czh",
		Params:      "{\"Prompt\": \"test_1\", \"Width\": 640}",
		Picture:     "test_17.png", //每次测试时修改一下，避免违反属性值唯一的约束
		Create_time: time.Now(),
	}
	db.DB.Create(&imageData)

	// 创建一个有效的Token
	claims := &middlewire.Claims{
		Username: "czh", //根据自己数据库已有的用户
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(), // 设置有效的过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(middlewire.JwtKey)

	req, _ := http.NewRequest("GET", "/auth/image/feature?feature=test_", nil)
	req.Header.Set("Authorization", tokenString)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotNil(t, response["images"]) // 确保有获取到数据
}

// TestFindByFeature_Success 测试成功(isOwn为true)获取图片的情况
func TestFindByFeature_Success2(t *testing.T) {
	fmt.Println()
	fmt.Println()
	fmt.Println("\n----------------------------------TestFindByFeature_Success")
	gin.SetMode(gin.TestMode)
	router := SetupRouter()

	// 插入测试数据
	imageData := image.ImageInformation{
		UserName:    "czh",
		Params:      "{\"Prompt\": \"test_1\", \"Width\": 640}",
		Picture:     "test_16.png", //每次测试时修改一下，避免违反属性值唯一的约束
		Create_time: time.Now(),
	}
	db.DB.Create(&imageData)

	// 创建一个有效的Token
	claims := &middlewire.Claims{
		Username: "czh", //根据自己数据库已有的用户
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(), // 设置有效的过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(middlewire.JwtKey)

	req, _ := http.NewRequest("GET", "/auth/image/feature?feature=test_&isOwn=true", nil)
	req.Header.Set("Authorization", tokenString)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotNil(t, response["images"]) // 确保有获取到数据
}

// TestFindByFeature_NoFeature 测试没有特征值的情况
func TestFindByFeature_NoFeature(t *testing.T) {
	fmt.Println()
	fmt.Println()
	fmt.Println("\n----------------------------------TestFindByFeature_NoFeature")
	gin.SetMode(gin.TestMode)
	router := SetupRouter()

	// 创建一个有效的Token
	claims := &middlewire.Claims{
		Username: "czh", //根据自己数据库已有的用户
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(), // 设置有效的过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(middlewire.JwtKey)

	req, _ := http.NewRequest("GET", "/auth/image/feature", nil)

	req.Header.Set("Authorization", tokenString)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, response["images"]) // 由于没有特征，返回值应为 nil
}

// TestFindByFeature_NoUserImages 测试用户没有相关图片的情况
func TestFindByFeature_NoUserImages(t *testing.T) {
	fmt.Println()
	fmt.Println()
	fmt.Println("\n----------------------------------TestFindByFeature_NoUserImages")
	gin.SetMode(gin.TestMode)
	router := SetupRouter()
	// 创建一个有效的Token
	claims := &middlewire.Claims{
		Username: "czh", //根据自己数据库已有的用户
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(), // 设置有效的过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(middlewire.JwtKey)

	req, _ := http.NewRequest("GET", "/auth/image/feature?feature=11111111111111111111&isOwn=true", nil) //模拟复杂特征
	req.Header.Set("Authorization", tokenString)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Len(t, response["images"], 0) // 由于用户没有图片，返回的图片列表应为空
}

// TestFindByFeature_InvalidToken 测试未找到用户信息的情况——即无效的token
func TestFindByFeature_InvalidToken(t *testing.T) {
	fmt.Println("\n----------------------------------InvalidToken")
	gin.SetMode(gin.TestMode)
	router := SetupRouter()
	// 创建一个无效的Token
	tokenString := "invalid-token"
	// 创建一个 GET 请求
	req, _ := http.NewRequest("GET", "/auth/image/feature?feature=test_&isOwn=true", nil)
	req.Header.Set("Authorization", tokenString)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 检查响应状态码
	assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected status code to be Unauthorized")

	// 解析响应体为map
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// 检查响应体中的'message'字段
	message, ok := response["message"].(string)
	if !ok {
		t.Error("expected 'message' to be a string")
		return
	}
	assert.Equal(t, "无效的Token", message, "Expected message to be '无效的Token'")
}

// TestFindByFeature_validToken 有效的token
func TestFindByFeature_validToken(t *testing.T) {
	fmt.Println("\n----------------------------------InvalidToken")
	gin.SetMode(gin.TestMode)
	router := SetupRouter()
	// 创建一个有效的Token
	claims := &middlewire.Claims{
		Username: "czh", //根据自己数据库已有的用户
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(), // 设置有效的过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(middlewire.JwtKey)

	req, _ := http.NewRequest("GET", "/auth/image/feature?feature=test_", nil) //模拟复杂特征
	req.Header.Set("Authorization", tokenString)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	fmt.Println("返回的状态码为：", w.Code)
	// 检查响应状态码不为401(401 表示“未找到用户信息”)
	assert.NotEqual(t, http.StatusUnauthorized, w.Code)
}

// TestFindByFeature_NoToken 测试缺少token的情况
func TestFindByFeature_NoToken(t *testing.T) {
	fmt.Println("\n----------------------------------NoToken")
	gin.SetMode(gin.TestMode)
	router := SetupRouter()

	// 创建一个 GET 请求
	req, _ := http.NewRequest("GET", "/auth/image/feature?feature=test_&isOwn=true", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 检查响应状态码
	assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected status code to be Unauthorized")

	// 检查响应体
	expectedResponse := gin.H{
		"code":    float64(401),
		"message": "请求头中缺少Token",
	}
	var actualResponse gin.H
	json.Unmarshal(w.Body.Bytes(), &actualResponse)
	assert.Equal(t, expectedResponse, actualResponse)
}
