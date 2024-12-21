package query

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
	"text-to-picture/models/image"
	db "text-to-picture/models/init"
	"text-to-picture/models/user"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

// SetupRouter2 设置 Gin 路由
func SetupRouter2() *gin.Engine {
	r := gin.Default()
	r.GET("/auth/user/images/timeRange", middlewire.JWTAuthMiddleware(), GetImagesWithinTimeRange)
	return r
}

// TestGetImagesWithinTimeRange 是测试 GetImagesWithinTimeRange 函数的入口函数
func TestGetImagesWithinTimeRange(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 读取测试数据库配置
	yamlFile, err := os.ReadFile(getDB.GetDBConfigPath())
	if err != nil {
		fmt.Printf("Error reading DBconfig.yaml file: %v\n", err)
		os.Exit(1)
	}

	// 解析配置文件
	var dbconfig DBConfig
	err = yaml.Unmarshal(yamlFile, &dbconfig)
	if err != nil {
		fmt.Printf("Error parsing DBconfig.yaml file: %v\n", err)
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

	// 清理数据库的函数
	defer func() {
		db.DB.Where("username = ?", "test").Delete(&user.UserInformation{})
		db.DB.Where("username = ?", "test").Delete(&image.ImageInformation{})
	}()

	// 创建路由
	router := SetupRouter2()

	// 测试缺少用户名
	t.Run("Missing_Username", func(t *testing.T) {

		router1 := gin.Default()
		// 模拟 JWT 中间件，但不设置用户名到上下文中
		router1.GET("/auth/user/images/timeRange", MockJWTAuthMiddlewareNoUser(), GetImagesWithinTimeRange)
	
		// 创建有效的Token
		claims := &middlewire.Claims{
			Username: "test",
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(middlewire.JwtKey)
	
		// 创建一个GET请求
		request, _ := http.NewRequest("GET", "/auth/user/images/timeRange?start_time=2023-10-01&end_time=2024-10-31", nil)
		request.Header.Set("Authorization", "Bearer "+tokenString)
	
		// 创建一个响应器
		response := httptest.NewRecorder()
	
		// 执行请求
		router1.ServeHTTP(response, request)
	
		// 检查响应码
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	
		// 检查响应体
		expectedResponse := gin.H{
			"code":    float64(401),
			"message": "未找到用户信息",
		}
		var actualResponse gin.H
		json.Unmarshal(response.Body.Bytes(), &actualResponse)
		fmt.Println("实际响应为：", actualResponse)
		assert.Equal(t, expectedResponse, actualResponse)
	})

	// 测试未找到图像
	t.Run("No_Images_Found", func(t *testing.T) {
		testUsername := "test"

		claims := &middlewire.Claims{
			Username: testUsername,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(middlewire.JwtKey)

		startTime := "2023-10-01T00:00:00Z"
		endTime := "2024-10-31T23:59:59Z"

		request, _ := http.NewRequest("GET", fmt.Sprintf("/auth/user/images/timeRange?start_time=%s&end_time=%s", startTime, endTime), nil)
		request.Header.Set("Authorization", tokenString)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)

		expectedResponse := gin.H{"code": float64(200), "message": "查询图像列表成功", "images": []interface{}{}}
		var actualResponse gin.H
		json.Unmarshal(response.Body.Bytes(), &actualResponse)
		fmt.Println("实际响应为：",actualResponse)
		assert.Equal(t, expectedResponse, actualResponse)
	})

	// 测试有效的时间范围，成功查询
	t.Run("Valid_Time_Range", func(t *testing.T) {
		testUsername := "test"

		claims := &middlewire.Claims{
			Username: testUsername,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(middlewire.JwtKey)

		// 插入测试数据
		imageInfo := &image.ImageInformation{
			UserName:   testUsername,
			Picture:    "http://example.com/test.jpg",
			Create_time: time.Now(),
		}
		db.DB.Create(&imageInfo)

		startTime := "2023-10-01T00:00:00Z"
		endTime := "2024-10-31T23:59:59Z"

		request, _ := http.NewRequest("GET", fmt.Sprintf("/auth/user/images/timeRange?start_time=%s&end_time=%s", startTime, endTime), nil)
		request.Header.Set("Authorization", tokenString)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)

		expectedResponse := gin.H{"code": float64(200), "message": "查询图像列表成功", "images": []interface{}{imageInfo}} 
		var actualResponse gin.H
		json.Unmarshal(response.Body.Bytes(), &actualResponse)
		fmt.Println("实际响应为：",actualResponse)
		assert.Equal(t, expectedResponse, actualResponse)

		// 清理测试数据
		db.DB.Where("username = ?", testUsername).Delete(&image.ImageInformation{})
	})
}
