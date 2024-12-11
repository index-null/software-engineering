package avator

//编写测试用例
import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"text-to-picture/middlewire/jwt"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"fmt"
	db "text-to-picture/models/init"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/auth/setavator",middlewire.JWTAuthMiddleware(),SetAvator)
	return r
}


func TestSetAvator(t *testing.T) {
	
	//写入DBConfig.yaml文件（数据库的配置文件）
	yamlFile , err := os.ReadFile("D:/go project/src/gocode/software-engineering/backend/text-to-picture/config/DBconfig/DBconfig.yaml")
	if err != nil {
		t.Fatalf("Error reading config.yaml file: %v", err)
		os.Exit(1)
	}

	//解析配置文件
	var dbconfig DBConfig
	err = yaml.Unmarshal(yamlFile, &dbconfig)
	if err != nil {
		t.Fatalf("Error parsing config.yaml file: %v", err)
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

	//测试 头像更新成功 的状态响应
	t.Run("SetAvator_Success", func(t *testing.T) {

		//设置gin为测试模式
		gin.SetMode(gin.TestMode)

		//创建一个路由
		router := setupRouter()

		// 创建一个有效的Token
		claims := &middlewire.Claims{
			Username: "testuser", // 根据自己数据库已有的用户
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(1 * time.Hour).Unix(), // 设置有效的过期时间
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(middlewire.JwtKey)
		
		//创建一个POST请求
		body := bytes.NewBuffer([]byte(`{"url": "https://example.com/new-avatar.jpg"}`))
		request , _ := http.NewRequest("POST", "/auth/setavator", body)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", tokenString)

		//创建一个响应器
		response := httptest.NewRecorder()

		//执行请求
		router.ServeHTTP(response, request)

		//检查响应码
		assert.Equal(t, http.StatusOK, response.Code)

		//检查响应体
		expectResponse := AvatorResponse{
			Code: Success,
			Msg:  "头像更新成功",
			Data: "https://example.com/new-avatar.jpg",
		}
		var actualResponse AvatorResponse
		json.Unmarshal(response.Body.Bytes(), &actualResponse)
		assert.Equal(t, expectResponse, actualResponse)
	})

	//测试 请求头缺失 的状态响应
	t.Run("SetAvator_MissingToken",func(t *testing.T) {

		//设置gin为测试模式
		gin.SetMode(gin.TestMode)

		//创建一个路由
		router := setupRouter()

		//创建一个POST请求
		body := bytes.NewBuffer([]byte(`{"url": "https://example.com/new-avatar.jpg"}`))
	    request , _ := http.NewRequest("POST", "/auth/setavator", body)
		request.Header.Set("Content-Type", "application/json")
		
		//创建一个响应器
		response := httptest.NewRecorder()

		//执行请求
		router.ServeHTTP(response, request)
	
		//检查响应码
		assert.Equal(t, Unauthorized , response.Code)
	
		//检查响应体
		expectResponse := gin.H{
			"code":    float64(http.StatusUnauthorized),
			"message": "请求头中缺少Token",
		}
		var actualResponse gin.H
		json.Unmarshal(response.Body.Bytes(), &actualResponse)
		assert.Equal(t, expectResponse, actualResponse)
	})

	//测试 无效的token 的状态响应
	t.Run("SetAvator_InvalidToken",func(t *testing.T) {

		//设置gin为测试模式
		gin.SetMode(gin.TestMode)

		//创建一个路由
		router := setupRouter()

		// 创建一个POST请求
		body := bytes.NewBuffer([]byte(`{"url": "https://example.com/new-avatar.jpg"}`))
	    request , _ := http.NewRequest("POST", "/auth/setavator", body)
		request.Header.Set("Content-Type", "application/json")

		//放入一个无效令牌
		request.Header.Set("Authorization", "Bearer invalid-token")

		//创建一个响应器
		response := httptest.NewRecorder()

		//执行请求
		router.ServeHTTP(response, request)

		//检查响应码
		assert.Equal(t, Unauthorized , response.Code)

		//检查响应体
		expectResponse := gin.H{
			"code":    float64(http.StatusUnauthorized),
			"message": "无效的Token",
		}
		var actualResponse gin.H
		json.Unmarshal(response.Body.Bytes(), &actualResponse)
		assert.Equal(t, expectResponse, actualResponse)
	})

	//更新头像失败
	t.Run("SetAvator_UpdateAvatorFailed",func(t *testing.T) {
		
		//设置gin为测试模式
		gin.SetMode(gin.TestMode)

		//创建一个路由
		router := setupRouter()
	
		// 创建一个有效的Token,但是该用户不存在
		claims := &middlewire.Claims{
			Username: "nonexistentuser", 
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(1 * time.Hour).Unix(), 
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(middlewire.JwtKey)
		
		//创建一个POST请求
		body := bytes.NewBuffer([]byte(`{"url": "https://example.com/new-avatar.jpg"}`))
	    request , _ := http.NewRequest("POST", "/auth/setavator", body)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", tokenString)

		//创建一个响应器
		response := httptest.NewRecorder()

		//执行请求
		router.ServeHTTP(response, request)

		//检查响应码
		assert.Equal(t, Error , response.Code)

		//检查响应体
		expectResponse := AvatorResponse{
			Code: Error,
			Msg:  "更新头像失败",
		}
		var actualResponse AvatorResponse
		json.Unmarshal(response.Body.Bytes(), &actualResponse)
		assert.Equal(t, expectResponse, actualResponse)
	})

	//名字解析出错的响应
	t.Run("SetAvator_NameParseError",func(t *testing.T) {
		
		//设置gin为测试模式
		gin.SetMode(gin.TestMode)

		//创建一个路由
		router := setupRouter()
	
		// 创建一个有效的Token,但Username = ""
		claims := &middlewire.Claims{
			Username: "", 
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(1 * time.Hour).Unix(), 
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(middlewire.JwtKey)
		
		//创建一个POST请求
		body := bytes.NewBuffer([]byte(`{"url": "https://example.com/new-avatar.jpg"}`))
	    request , _ := http.NewRequest("POST", "/auth/setavator", body)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", tokenString)

		//创建一个响应器
		response := httptest.NewRecorder()

		//执行请求
		router.ServeHTTP(response, request)

		//检查响应码
		assert.Equal(t, http.StatusUnauthorized, response.Code)

		//检查响应体
		expectResponse := AvatorResponse{
			Code: Unauthorized,
			Msg:  "名字解析出错",
		}
		var actualResponse AvatorResponse
		json.Unmarshal(response.Body.Bytes(), &actualResponse)
		assert.Equal(t, expectResponse, actualResponse)
	})
}


