package like

import(
	"testing"
	"text-to-picture/models/image"
	"github.com/gin-gonic/gin"
	db "text-to-picture/models/init"
	"os"
	"fmt"
	"gopkg.in/yaml.v3"
	"encoding/json"
	"bytes"
	"net/http"
	
	"text-to-picture/middlewire/jwt"
	"time"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)


func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/auth/like",middlewire.JWTAuthMiddleware(),LikeImage)
	return r
}

// MockUserRepository 是一个模拟的图片点赞仓库
type MockImageLikeRepository struct{}

// GetUserByName 是一个模拟的获取用户方法
func (m *MockImageLikeRepository) GetUserByName(name string) (*image.ImageLike, error) {
	// 返回一个模拟的用户信息
	return &image.ImageLike{
		UserName: "testuser",
		Picture:  "generate/string-2024-12-09 15:36:59.png",
		Num:      0,
	}, nil
}

type mockDB struct {
    *gorm.DB
    beginError error
}

// 实现 gorm.DB 的 Begin 方法
func (m *mockDB) Begin() *gorm.DB {
    if m.beginError != nil {
        return nil
    }
    return m.DB.Begin()
}

// 创建一个点赞
func createLikeRequest(url string, tokenString string) (*http.Request, *bytes.Buffer) {
    body := bytes.NewBuffer([]byte(fmt.Sprintf(`{"url": "%s"}`, url)))
    request, err := http.NewRequest("POST", "/auth/like", body)
    if err != nil {
        panic(err) // 或者根据实际情况处理错误
    }
    request.Header.Set("Content-Type", "application/json")
    request.Header.Set("Authorization", tokenString)
	//创建一个响应器
	response := httptest.NewRecorder()

	//执行请求
	r := SetupRouter()
	r.ServeHTTP(response, request)
    return request, body
}

// 定义一个全局变量来存储 getLikeCount 的引用
var getLikeCountFunc func(string) int

// 初始化全局变量
func init() {
    getLikeCountFunc = getLikeCount
}

//获取点赞数
func getLikeCount(url string) int {
	// 获取当前点赞数
	var currentLikeCount int
	if err := db.DB.Model(&image.ImageInformation{}).Where("picture = ?", url).Select("likecount").Row().Scan(&currentLikeCount); err != nil {
		return 0
	}
	return currentLikeCount
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
func TestLike(t *testing.T) {

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

	//要点赞的图片的url,这里可以自己去生成一张图片，获取它的url
	var imageUrl = "generate/string-2024-12-09 20:12:59.png"

	//成功的响应
	t.Run("LikeImage_Success", func(t *testing.T) {

		//设置gin为测试模式
		gin.SetMode(gin.TestMode)

		//创建一个路由
		router := SetupRouter()

		// 创建一个有效的Token
		claims := &middlewire.Claims{
			Username: "testuser", // 根据自己数据库已有的且为点赞该图片的用户
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(1 * time.Hour).Unix(), // 设置有效的过期时间
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(middlewire.JwtKey)
	
		//创建一个POST请求
		body := bytes.NewBuffer([]byte(fmt.Sprintf(`{"url": "%s"}`, imageUrl)))
		request , _ := http.NewRequest("POST", "/auth/like", body)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", tokenString)

		//创建一个响应器
		response := httptest.NewRecorder()

		//获取该图片当前的点赞数
		currentLikeCount := getLikeCount(imageUrl)

		//执行请求
		router.ServeHTTP(response, request)

		//检查响应码
		assert.Equal(t, http.StatusOK, response.Code)

		//检查响应体
		expectResponse := gin.H{
			"current_likes": (float64)(currentLikeCount + 1 ),
		    "message":"Image liked successfully",
		}
		fmt.Println(currentLikeCount)
		var actualResponse gin.H
		json.Unmarshal(response.Body.Bytes(), &actualResponse)
		assert.Equal(t, expectResponse, actualResponse)
	})

	//缺少图片url的响应
	t.Run("LikeImage_MissImageURL", func(t *testing.T) {
	
		//设置gin为测试模式
		gin.SetMode(gin.TestMode)

		//创建一个路由
		router := SetupRouter()

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
		body := bytes.NewBuffer([]byte(`{"url": ""}`))
		request , _ := http.NewRequest("POST", "/auth/like", body)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", tokenString)

		//创建一个响应器
		response := httptest.NewRecorder()

		//执行请求
		router.ServeHTTP(response, request)

		//检查响应码
		assert.Equal(t, http.StatusBadRequest, response.Code)

		//检查响应体
		expectResponse := gin.H{
			"code":  float64(http.StatusBadRequest),
			"error": "Missing image URL",
		}
		var actualResponse gin.H
		json.Unmarshal(response.Body.Bytes(), &actualResponse)
		assert.Equal(t, expectResponse, actualResponse)
	})

	//名字解析错误的响应
	t.Run("LikeImage_NameParseError", func(t *testing.T) {
	
		//设置gin为测试模式
		gin.SetMode(gin.TestMode)

		//创建一个路由
		router := SetupRouter()

		// 创建一个有效的Token
		claims := &middlewire.Claims{
			Username: "", 
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(1 * time.Hour).Unix(), // 设置有效的过期时间
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(middlewire.JwtKey)
		
		//创建一个POST请求
		body := bytes.NewBuffer([]byte(fmt.Sprintf(`{"url": "%s"}`, imageUrl)))
		request , _ := http.NewRequest("POST", "/auth/like", body)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", tokenString)

		//创建一个响应器
		response := httptest.NewRecorder()

		//执行请求
		router.ServeHTTP(response, request)

		//检查响应码
		assert.Equal(t, http.StatusUnauthorized , response.Code)

		//检查响应体
		expectResponse := gin.H{
			"code":  float64(http.StatusUnauthorized ),
			"error": "名字解析出错",
		}
		var actualResponse gin.H
		json.Unmarshal(response.Body.Bytes(), &actualResponse)
		assert.Equal(t, expectResponse, actualResponse)
	})

	//点赞数据库开始出错
	t.Run("LikeImage_DBBeginError", func(t *testing.T) {
	
		//设置gin为测试模式
		gin.SetMode(gin.TestMode)

		//创建一个路由
		router := SetupRouter()

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
		body := bytes.NewBuffer([]byte(fmt.Sprintf(`{"url": "%s"}`, imageUrl)))
		request , _ := http.NewRequest("POST", "/auth/like", body)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", tokenString)

		//创建一个响应器
		response := httptest.NewRecorder()

		originalDB := db.DB
		//db.DB = &mockDB{DB: db.DB, beginError: fmt.Errorf("database begin error")}

		//执行请求
		router.ServeHTTP(response, request)

		db.DB = originalDB

		//检查响应码
		assert.Equal(t, http.StatusInternalServerError, response.Code)

		//检查响应体
		expectResponse := gin.H{
			"code":  float64(http.StatusInternalServerError),
			"error": "点赞数据库开始出错",
		}
		var actualResponse gin.H
		json.Unmarshal(response.Body.Bytes(), &actualResponse)
		assert.Equal(t, expectResponse, actualResponse)
	})

	//用户已经点赞过的响应
t.Run("LikeImage_AlreadyLiked", func(t *testing.T) {
	
	//设置gin为测试模式
	gin.SetMode(gin.TestMode)

	//创建一个路由
	router := SetupRouter()

	// 创建一个有效的Token
	claims := &middlewire.Claims{
		Username: "testuser2", // 根据自己数据库已有的用户
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(), // 设置有效的过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(middlewire.JwtKey)
	
	//创建一个点赞，使该用户再次点赞时会出现已经点赞过的错误
	createLikeRequest(imageUrl,tokenString)

	//创建一个POST请求
	body := bytes.NewBuffer([]byte(fmt.Sprintf(`{"url": "%s"}`, imageUrl)))
	request , _ := http.NewRequest("POST", "/auth/like", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", tokenString)

	//创建一个响应器
	response := httptest.NewRecorder()

	//执行请求
	router.ServeHTTP(response, request)

	//检查响应码
	assert.Equal(t, http.StatusConflict, response.Code)

	//检查响应体
	expectResponse := gin.H{
		"code":  float64(http.StatusConflict),
		"error": "用户已经点赞过该图片",
	}
	var actualResponse gin.H
	json.Unmarshal(response.Body.Bytes(), &actualResponse)
	assert.Equal(t, expectResponse, actualResponse)
})

//获取点赞失败的响应
t.Run("LikeImage_GetLikeCountFault", func(t *testing.T) {
	
	//设置gin为测试模式
	gin.SetMode(gin.TestMode)

	//创建一个路由
	router := SetupRouter()

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
	body := bytes.NewBuffer([]byte(`{"url": "test.png"}`))//这里就放一个错误的url
	request , _ := http.NewRequest("POST", "/auth/like", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", tokenString)

	//创建一个响应器
	response := httptest.NewRecorder()

	 // 模拟 getLikeCount 函数返回错误
	 originalGetLikeCount := getLikeCountFunc
	 getLikeCountFunc = func(url string) int {
		 return -1 // 返回一个无效值以模拟错误
	 }
	 defer func() {
		 getLikeCountFunc = originalGetLikeCount // 测试结束后恢复原始函数
	 }()

	//执行请求
	router.ServeHTTP(response, request)

	//检查响应码
	assert.Equal(t, http.StatusInternalServerError, response.Code)

	//检查响应体
	expectResponse := gin.H{
		"code":  float64(http.StatusInternalServerError),
		"error": "sql: no rows in result set",
	}
	var actualResponse gin.H
	json.Unmarshal(response.Body.Bytes(), &actualResponse)
	assert.Equal(t, expectResponse, actualResponse)
})


}




