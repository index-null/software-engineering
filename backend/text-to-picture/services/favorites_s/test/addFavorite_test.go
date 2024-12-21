package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	//"os"
	"testing"

	getDB "text-to-picture/config"
	middlewire "text-to-picture/middlewire/jwt"
	"text-to-picture/models/image"
	db "text-to-picture/models/init"
	"text-to-picture/models/user"
	"text-to-picture/services/favorites_s"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

// SetupRouter 设置 Gin 路由
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/auth/addFavoritedImage", middlewire.JWTAuthMiddleware(), favorites_s.AddFavoritedImage)
	return r
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

// TestAddFavoritedImage 是测试 AddFavoritedImage 函数的入口函数
func TestAddFavoritedImage(t *testing.T) {
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
		// 清理收藏的图像
		db.DB.Where("username = ?", "test").Delete(&image.FavoritedImages{}) // 假设有一个 FavoritedImage 模型
	}()

	// 创建路由
	router := SetupRouter()

	// 测试无效的请求格式
	t.Run("Invalid_Request_Format", func(t *testing.T) {
		// 创建一个数据库尚未存在的用户作为测试用户
		testUsername := "test" 

		// 创建有效的Token
		claims := &middlewire.Claims{
			Username: testUsername, 
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(middlewire.JwtKey)
		//向数据库插入该用户的信息
		create_time := time.Now().UTC()
		db.DB.Create(&user.UserInformation{ID:100,UserName: testUsername, Email:testUsername+"@qq.com",Password: "123456",Avatar_url:testUsername+".jpg",Score: 100,Token: tokenString,Create_time: create_time,})

		request, _ := http.NewRequest("POST", "/auth/addFavoritedImage", nil)
		request.Header.Set("Authorization", tokenString)

		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)

		expectedResponse := gin.H{"message": "无效的请求格式", "error": "invalid request"}
		var actualResponse gin.H
		json.Unmarshal(response.Body.Bytes(), &actualResponse)
		fmt.Println("实际响应：", actualResponse)
		assert.Equal(t, expectedResponse, actualResponse)

		//最后清理测试数据
		db.DB.Where("username = ?", testUsername).Delete(&user.UserInformation{})
	})

	// 测试无效的图像 URL 或 ID
	t.Run("Invalid_Image_URL_Or_ID", func(t *testing.T) {
		// 创建一个数据库尚未存在的用户作为测试用户
		testUsername := "test" 

		// 创建有效的Token
		claims := &middlewire.Claims{
			Username: testUsername, 
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(middlewire.JwtKey)
		//向数据库插入该用户的信息
		create_time := time.Now().UTC()
		db.DB.Create(&user.UserInformation{ID:100,UserName: testUsername, Email:testUsername+"@qq.com",Password: "123456",Avatar_url:testUsername+".jpg",Score: 100,Token: tokenString,Create_time: create_time,})

		requestBody := `{"url": "", "id": 0}`
		request, _ := http.NewRequest("POST", "/auth/addFavoritedImage", strings.NewReader(requestBody))
		request.Header.Set("Authorization", tokenString)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)

		expectedResponse := gin.H{"message": "无有效的图像id或url", "error": "id 必须大于 0 或者 url 不得为空"}
		var actualResponse gin.H
		json.Unmarshal(response.Body.Bytes(), &actualResponse)
		fmt.Println("实际响应：", actualResponse)
		assert.Equal(t, expectedResponse, actualResponse)

		//最后清理测试数据
		db.DB.Where("username = ?", testUsername).Delete(&user.UserInformation{})
	})

	// 测试未找到对应的图像
	t.Run("Image_Not_Found", func(t *testing.T) {
		// 创建一个数据库尚未存在的用户作为测试用户
		testUsername := "test" 

		// 创建有效的Token
		claims := &middlewire.Claims{
			Username: testUsername, 
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(middlewire.JwtKey)
		//向数据库插入该用户的信息
		create_time := time.Now().UTC()
		db.DB.Create(&user.UserInformation{ID:100,UserName: testUsername, Email:testUsername+"@qq.com",Password: "123456",Avatar_url:testUsername+".jpg",Score: 100,Token: tokenString,Create_time: create_time,})

		requestBody := `{"url": "http://example.com/notfound.jpg", "id": 0}`
		request, _ := http.NewRequest("POST", "/auth/addFavoritedImage", strings.NewReader(requestBody))
		request.Header.Set("Authorization", tokenString)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNotFound, response.Code)

		expectedResponse := gin.H{"message": "未找到对应的图像", "error": "record not found"}
		var actualResponse gin.H
		json.Unmarshal(response.Body.Bytes(), &actualResponse)
		fmt.Println("实际响应：", actualResponse)
		assert.Equal(t, expectedResponse, actualResponse)

		//最后清理测试数据
		db.DB.Where("username = ?", testUsername).Delete(&user.UserInformation{})
	})

	// 测试该图像已经被收藏过
	t.Run("Image_Already_Favorited", func(t *testing.T) {
		// 创建一个数据库尚未存在的用户作为测试用户
		testUsername := "test" 

		// 创建有效的Token
		claims := &middlewire.Claims{
			Username: testUsername, 
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(middlewire.JwtKey)
		//向数据库插入该用户的信息
		create_time := time.Now().UTC()
		db.DB.Create(&user.UserInformation{ID:100,UserName: testUsername, Email:testUsername+"@qq.com",Password: "123456",Avatar_url:testUsername+".jpg",Score: 100,Token: tokenString,Create_time: create_time,})

		// 先插入测试数据
		imageUrl := "http://example.com/test.jpg"
		imageInfo := &image.ImageInformation{
			UserName: testUsername,
			Picture:     imageUrl,
			Create_time: time.Now(),
		}
		db.DB.Create(&imageInfo)

		// 假设用户已收藏该图像
		db.DB.Create(&image.FavoritedImages{UserName: testUsername, Picture: imageUrl})

		requestBody := `{"url": "http://example.com/test.jpg", "id": 0}`
		request, _ := http.NewRequest("POST", "/auth/addFavoritedImage", strings.NewReader(requestBody))
		request.Header.Set("Authorization", tokenString)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusConflict, response.Code)

		expectedResponse := gin.H{"message": "该图像已经被收藏过"}
		var actualResponse gin.H
		json.Unmarshal(response.Body.Bytes(), &actualResponse)
		fmt.Println("实际响应：", actualResponse)
		assert.Equal(t, expectedResponse, actualResponse)

		//最后清理测试数据
		db.DB.Where("username = ?", testUsername).Delete(&user.UserInformation{})
		db.DB.Where("username = ?", testUsername).Delete(&image.ImageInformation{})
		db.DB.Where("username = ?", testUsername).Delete(&image.FavoritedImages{})
	})

	// 测试收藏成功
	t.Run("Add_Favorited_Image_Success", func(t *testing.T) {
		// 创建一个数据库尚未存在的用户作为测试用户
		testUsername := "test" 

		// 创建有效的Token
		claims := &middlewire.Claims{
			Username: testUsername, 
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(middlewire.JwtKey)
		//向数据库插入该用户的信息
		create_time := time.Now().UTC()
		db.DB.Create(&user.UserInformation{ID:100,UserName: testUsername, Email:testUsername+"@qq.com",Password: "123456",Avatar_url:testUsername+".jpg",Score: 100,Token: tokenString,Create_time: create_time,})

		// 先插入测试数据
		imageUrl := "http://example.com/newimage.jpg"
		imageInfo := &image.ImageInformation{
			UserName: testUsername,
			Picture:     imageUrl,
			Create_time: time.Now(),
		}
		db.DB.Create(&imageInfo)

		requestBody := `{"url": "http://example.com/newimage.jpg", "id": 0}`
		request, _ := http.NewRequest("POST", "/auth/addFavoritedImage", strings.NewReader(requestBody))
		request.Header.Set("Authorization", tokenString)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)

		expectedResponse := gin.H{"message": "图像收藏成功"}
		var actualResponse gin.H
		json.Unmarshal(response.Body.Bytes(), &actualResponse)
		fmt.Println("实际响应：", actualResponse)
		assert.Equal(t, expectedResponse, actualResponse)

		//最后清理测试数据
		db.DB.Where("username = ?", testUsername).Delete(&user.UserInformation{})
		db.DB.Where("username = ?", testUsername).Delete(&image.ImageInformation{})
		db.DB.Where("username = ?", testUsername).Delete(&image.FavoritedImages{})
	})
}
