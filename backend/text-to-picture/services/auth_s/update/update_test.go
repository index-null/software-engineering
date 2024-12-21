package update

import (
	"bytes"
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
	"text-to-picture/models/user"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

type DBConfig struct {
	DB struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"db"`
}

// SetupRouter 设置 Gin 路由
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.PUT("/auth/user/update", middlewire.JWTAuthMiddleware(), UpdateUser)
	return r
}
// MockJWTAuthMiddlewareNoUser 模拟 JWT 中间件，但不设置用户名到上下文中
func MockJWTAuthMiddlewareNoUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func TestUpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 读取测试数据库配置
	yamlFile, err := os.ReadFile(getDB.GetDBConfigPath())
	if err != nil {
		t.Fatalf("Error reading DBconfig.yaml file: %v", err)
	}

	// 解析配置文件
	var dbconfig DBConfig
	err = yaml.Unmarshal(yamlFile, &dbconfig)
	if err != nil {
		t.Fatalf("Error parsing DBconfig.yaml file: %v", err)
	}

	// 设置数据库连接的环境变量
	os.Setenv("DB_USER", dbconfig.DB.User)
	os.Setenv("DB_PASSWORD", dbconfig.DB.Password)
	os.Setenv("DB_HOST", dbconfig.DB.Host)
	os.Setenv("DB_PORT", dbconfig.DB.Port)
	os.Setenv("DB_NAME", dbconfig.DB.Name)

	// 连接数据库
	if err := db.ConnectDatabase(); err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// 创建路由
	router := SetupRouter()

	testUsername := "test_update"

	// 清理数据库的函数
	defer func() {
		db.DB.Where("username = ?", testUsername).Delete(&user.UserInformation{})
	}()

	// 插入测试数据
	userinfo := &user.UserInformation{
		UserName: testUsername,
		Password: "123456",
		Email:    testUsername + "@qq.com",
	}
	db.DB.Create(&userinfo)

	claims := &middlewire.Claims{
		Username: testUsername,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(middlewire.JwtKey)

	// 辅助函数来创建请求并执行测试
	executeRequest := func(requestBody map[string]interface{}, t *testing.T) (*httptest.ResponseRecorder, error) {
		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest("PUT", "/auth/user/update", bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", tokenString)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		return resp, nil
	}

	// 测试成功更新用户信息
	t.Run("Successful_Update", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"password": "abcdefg",
			"email":    "newemail@qq.com",
			"avatar_url": "http://example.com/new_avatar.jpg",
		}

		resp, err := executeRequest(requestBody, t)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		var actualResponse gin.H
		if err := json.Unmarshal(resp.Body.Bytes(), &actualResponse); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}
		fmt.Println("实际响应为：", actualResponse)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "用户信息更新成功", actualResponse["message"])
	})

	// 测试尝试更新用户名（应失败）
	t.Run("AttemptToUpdateUsername", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"username": "new_username",
		}

		resp, err := executeRequest(requestBody, t)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		var actualResponse gin.H
		if err := json.Unmarshal(resp.Body.Bytes(), &actualResponse); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}
		fmt.Println("实际响应为：", actualResponse)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Contains(t, actualResponse["error"], "用户名不可修改")
	})

	// 测试密码长度小于6位
	t.Run("PasswordTooShort", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"password": "123",
		}

		resp, err := executeRequest(requestBody, t)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		var actualResponse gin.H
		if err := json.Unmarshal(resp.Body.Bytes(), &actualResponse); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}
		fmt.Println("用户名，实际响应为：", actualResponse)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Contains(t, actualResponse["error"], "密码少于6位")
	})

	// 测试邮箱格式错误
	t.Run("InvalidEmailFormat", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"email": "123qq.com",
		}

		resp, err := executeRequest(requestBody, t)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		var actualResponse gin.H
		if err := json.Unmarshal(resp.Body.Bytes(), &actualResponse); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}
		fmt.Println("邮箱，实际响应为：", actualResponse)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Contains(t, actualResponse["error"], "邮箱格式不正确")
	})

	// 测试邮箱为空
	t.Run("EmptyEmail", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"email": "",
		}

		resp, err := executeRequest(requestBody, t)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		var actualResponse gin.H
		if err := json.Unmarshal(resp.Body.Bytes(), &actualResponse); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}
		fmt.Println("邮箱，实际响应为：", actualResponse)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Contains(t, actualResponse["error"], "邮箱为空")
	})

	// 测试无效的请求体
	t.Run("Invalid_Request_Body", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/auth/user/update", bytes.NewBuffer([]byte(`invalid json`)))
		req.Header.Set("Authorization", tokenString)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		var actualResponse gin.H
		if err := json.Unmarshal(resp.Body.Bytes(), &actualResponse); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}
		fmt.Println("实际响应为：", actualResponse)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, actualResponse["message"], "请求数据格式错误")
	})

	// 测试缺少认证信息
	t.Run("Missing_Authentication", func(t *testing.T) {
		router1 := gin.Default()
		router1.PUT("/auth/user/update", MockJWTAuthMiddlewareNoUser(), UpdateUser)
		requestBody := map[string]interface{}{
			"email": "newemail@qq.com",
		}

		jsonData, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("PUT", "/auth/user/update", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", tokenString)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router1.ServeHTTP(resp, req)

		var actualResponse gin.H
		if err := json.Unmarshal(resp.Body.Bytes(), &actualResponse); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}
		fmt.Println("实际响应为：", actualResponse)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.Contains(t, actualResponse["message"], "未找到用户信息")
	})
}