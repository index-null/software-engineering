package delete

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	getDB "text-to-picture/config"
	middlewire "text-to-picture/middlewire/jwt"
	db "text-to-picture/models/init"
	u "text-to-picture/models/user"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
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
	r.DELETE("/auth/root/deleteOneUser", middlewire.JWTAuthMiddleware(), DeleteUserByName)
	return r
}

// TestDeleteUserByName_Success 测试root成功删除用户的情况
func TestDeleteUserByName_Success(t *testing.T) {
	fmt.Println("\n--------------------------------------------------------TestDeleteUserByName_Success")
	gin.SetMode(gin.TestMode)
	router := SetupRouter()

	// 创建测试用户
	testUsername := "testuser_success"
	db.DB.Create(&u.UserInformation{UserName: testUsername, Email: "testuser_success@qq.com"})

	// 创建有效的Token
	claims := &middlewire.Claims{
		Username: "root", // 使用root用户进行删除
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(middlewire.JwtKey)

	req, _ := http.NewRequest("DELETE", "/auth/root/deleteOneUser?username=testuser_success", nil)
	req.Header.Set("Authorization", tokenString)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "成功删除用户：testuser_success", response["message"])
}

// TestDeleteUserByName_Success 测试用户账号成功注销的情况
func TestDeleteUserByName_LogoutSuccess(t *testing.T) {
	fmt.Println("\n--------------------------------------------------------TestDeleteUserByName_LogoutSuccess")
	gin.SetMode(gin.TestMode)
	router := SetupRouter()

	// 创建测试用户
	testUsername := "testuser_logout"
	db.DB.Create(&u.UserInformation{UserName: testUsername, Email: "testuser_logout@qq.com"})

	// 创建有效的Token
	claims := &middlewire.Claims{
		Username: "testuser_logout", // 使用用户自身进行账号注销
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(middlewire.JwtKey)

	req, _ := http.NewRequest("DELETE", "/auth/root/deleteOneUser?username=testuser_logout&isOwn=true", nil)
	req.Header.Set("Authorization", tokenString)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	fmt.Println("返回的message为：", response["message"])
	assert.Equal(t, "testuser_logout的账号注销成功", response["message"])
}

// TestDeleteUserByName_UserNotFound 测试用户不存在的情况
func TestDeleteUserByName_UserNotFound(t *testing.T) {
	fmt.Println("\n--------------------------------------------------------TestDeleteUserByName_UserNotFound")
	gin.SetMode(gin.TestMode)
	router := SetupRouter()

	// 创建有效的Token
	claims := &middlewire.Claims{
		Username: "root", // 使用root用户进行删除
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(middlewire.JwtKey)

	req, _ := http.NewRequest("DELETE", "/auth/root/deleteOneUser?username=nonexistentuser", nil)
	req.Header.Set("Authorization", tokenString)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "用户不存在", response["message"])
}

// TestDeleteUserByName_NoPermission 测试非root用户尝试删除其他用户的情况
func TestDeleteUserByName_NoPermission(t *testing.T) {
	fmt.Println("\n--------------------------------------------------------TestDeleteUserByName_NoPermission")
	gin.SetMode(gin.TestMode)
	router := SetupRouter()

	// 创建测试用户
	testUsername := "testuser_non3" // 每次测试要修改，避免违反唯一约束"userinformation_username_key" 和"userinformation_email_key"
	db.DB.Create(&u.UserInformation{UserName: testUsername, Email: "testuser_non3@qq.com"})

	// 创建有效的Token
	claims := &middlewire.Claims{
		Username: "nonrootuser", // 非root用户
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(middlewire.JwtKey)

	req, _ := http.NewRequest("DELETE", "/auth/root/deleteOneUser?username=testuser_non3", nil)
	req.Header.Set("Authorization", tokenString)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "非root用户，不可删除其他某个用户", response["message"])
}

// TestDeleteUserByName_NoToken 测试缺少token的情况
func TestDeleteUserByName_NoToken(t *testing.T) {
	fmt.Println("\n--------------------------------------------------------TestDeleteUserByName_NoToken")
	gin.SetMode(gin.TestMode)
	router := SetupRouter()

	// 创建一个 DELETE 请求
	req, _ := http.NewRequest("DELETE", "/auth/root/deleteOneUser?username=notoken", nil)

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

// TestDeleteUserByName_InvalidToken 测试无效Token的情况
func TestDeleteUserByName_InvalidToken(t *testing.T) {
	fmt.Println("\n--------------------------------------------------------TestDeleteUserByName_InvalidToken")
	gin.SetMode(gin.TestMode)
	router := SetupRouter()

	req, _ := http.NewRequest("DELETE", "/auth/root/deleteOneUser?username=testuser", nil)
	req.Header.Set("Authorization", "invalid-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "无效的Token", response["message"])
}
