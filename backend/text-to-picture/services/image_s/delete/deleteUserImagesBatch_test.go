package delete

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
	"text-to-picture/models/image"
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
	r.POST("/auth/user/deleteImages", middlewire.JWTAuthMiddleware(), DeleteUserImagesBatch)
	return r
}

// MockJWTAuthMiddlewareNoUser 模拟 JWT 中间件，但不设置用户名到上下文中
func MockJWTAuthMiddlewareNoUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// TestDeleteUserImagesBatch 测试 DeleteUserImagesBatch 函数
func TestDeleteUserImagesBatch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 读取测试数据库配置
	yamlFile, err := os.ReadFile(getDB.GetDBConfigPath())
	if err != nil {
		t.Fatalf("Error reading configs.yaml file: %v", err)
	}

	// 解析配置文件
	var dbconfig DBConfig
	err = yaml.Unmarshal(yamlFile, &dbconfig)
	if err != nil {
		t.Fatalf("Error parsing configs.yaml file: %v", err)
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

	testUsername := "test_user"

	// 清理数据库的函数
	defer func() {
		db.DB.Where("username = ?", testUsername).Delete(&user.UserInformation{})
		db.DB.Where("username = ?", testUsername).Delete(&image.ImageInformation{})
	}()

	claims := &middlewire.Claims{
		Username: testUsername,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(middlewire.JwtKey)

	// 插入测试数据
	userinfo := &user.UserInformation{
		UserName: testUsername,
		Password: "123456",
		Email:    testUsername + "@qq.com",
	}
	imageInfo1 := &image.ImageInformation{
		UserName:    testUsername,
		Picture:     "http://example.com/test_1.jpg",
		Create_time: time.Now(),
	}
	imageInfo2 := &image.ImageInformation{
		UserName:    testUsername,
		Picture:     "http://example.com/test_2.jpg",
		Create_time: time.Now(),
	}
	db.DB.Create(&userinfo)
	db.DB.Create(&imageInfo1)
	db.DB.Create(&imageInfo2)

	// 测试成功删除
	t.Run("Successful_Delete", func(t *testing.T) {
		requestBody := BatchDeleteRequestBody{
			Ids: []int{int(imageInfo2.ID)},
		}

		jsonData, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/auth/user/deleteImages", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", tokenString)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		expectedResponse := gin.H{
			"message": "成功删除用户指定的图像",
		}
		var actualResponse gin.H
		json.Unmarshal(resp.Body.Bytes(), &actualResponse)
		fmt.Println("实际响应为：", actualResponse)
		assert.Equal(t, expectedResponse, actualResponse)

		requestBody = BatchDeleteRequestBody{
			Urls: []string{imageInfo1.Picture},
		}

		jsonData, _ = json.Marshal(requestBody)
		req, _ = http.NewRequest("POST", "/auth/user/deleteImages", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", tokenString)
		req.Header.Set("Content-Type", "application/json")
		resp = httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		expectedResponse = gin.H{
			"message": "成功删除用户指定的图像",
		}
		json.Unmarshal(resp.Body.Bytes(), &actualResponse)
		fmt.Println("实际响应为：", actualResponse)
		assert.Equal(t, expectedResponse, actualResponse)
	})

	// 测试无效的请求体
	t.Run("Invalid_Request_Body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/auth/user/deleteImages", bytes.NewBuffer([]byte(`invalid json`)))
		req.Header.Set("Authorization", tokenString)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		var actualResponse gin.H
		if err := json.Unmarshal(resp.Body.Bytes(), &actualResponse); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		fmt.Println("实际响应为：", actualResponse)

		// 检查 "message" 字段是否符合预期
		assert.Equal(t, "无效的请求格式", actualResponse["message"], "Expected message '无效的请求格式'")

		// 检查 "error" 字段是否存在且非空
		errorValue, exists := actualResponse["error"].(string)
		assert.True(t, exists, "Expected 'error' field to exist")
		assert.NotEmpty(t, errorValue, "Expected non-empty error field")
	})

	// 测试缺少用户名
	t.Run("Missing_Username", func(t *testing.T) {
		router1 := gin.Default()
		router1.POST("/auth/user/deleteImages", MockJWTAuthMiddlewareNoUser(), DeleteUserImagesBatch)

		requestBody := BatchDeleteRequestBody{
			Urls: []string{"http://example.com/test_1.jpg"},
		}

		jsonData, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/auth/user/deleteImages", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", tokenString)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router1.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)

		expectedResponse := gin.H{
			"success": false,
			"message": "未找到用户信息",
		}
		var actualResponse gin.H
		json.Unmarshal(resp.Body.Bytes(), &actualResponse)
		fmt.Println("实际响应为：", actualResponse)
		assert.Equal(t, expectedResponse, actualResponse)
	})

	// 测试同时提供 ids 和 urls
	t.Run("Both_Ids_and_URLs_Provided", func(t *testing.T) {
		requestBody := BatchDeleteRequestBody{
			Urls: []string{imageInfo1.Picture},
			Ids:  []int{imageInfo2.ID},
		}

		jsonData, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/auth/user/deleteImages", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", tokenString)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		expectedResponse := gin.H{
			"message": "请提供有效的urls或ids列表，并且不要同时提供这两个列表",
		}
		var actualResponse gin.H
		json.Unmarshal(resp.Body.Bytes(), &actualResponse)
		fmt.Println("实际响应为：", actualResponse)
		assert.Equal(t, expectedResponse, actualResponse)
	})
}