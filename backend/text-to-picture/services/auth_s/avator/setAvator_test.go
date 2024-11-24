package avator

//编写测试用例

import (
	"bytes"
	"encoding/json"

	"net/http"
	"net/http/httptest"
	"testing"
	"text-to-picture/middlewire/jwt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

/* //生成token
func generateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &middlewire.Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middlewire.JwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
} */

//设置路由
func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/auth/setavator", SetAvator)
	return r
}



//错误的更换头像的用例

//请求头缺少token

func TestSetAvator_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := setupRouter()
	
	body := bytes.NewBuffer([]byte(`{"url":"https://example.com/new-avatar.jpg"}`))
	request , _ := http.NewRequest("POST", "/auth/setavator", body)
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, Unauthorized , response.Code)

	expectResponse := AvatorResponse{
		Code: Unauthorized,
		Msg:  "请求头中缺少Token",
	}
	var actualResponse AvatorResponse
	json.Unmarshal(response.Body.Bytes(), &actualResponse)
	assert.Equal(t, expectResponse, actualResponse)
}


//无效token

func TestSetAvator_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := setupRouter()

	body := bytes.NewBuffer([]byte(`{"url":"https://example.com/new-avatar.jpg"}`))
	request , _ := http.NewRequest("POST", "/auth/setavator", body)
	request.Header.Set("Content-Type", "application/json")
	//放入一个无效令牌
	request.Header.Set("Authorization", "Bearer invalid-token")

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, Unauthorized , response.Code)

	expectResponse := AvatorResponse{
		Code: Unauthorized,
		Msg:  "无效的Token",
	}
	var actualResponse AvatorResponse
	json.Unmarshal(response.Body.Bytes(), &actualResponse)
	assert.Equal(t, expectResponse, actualResponse)
}

//更新头像失败
func TestSetAvator_UpdateAvatorFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := setupRouter()

	/* token , _ := generateToken("testuser") */
	// 创建一个有效的Token
	claims := &middlewire.Claims{
		Username: "qin", // 根据自己数据库已有的用户
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(), // 设置有效的过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(middlewire.JwtKey)
	
	body := bytes.NewBuffer([]byte(`{"url":"https://example.com/new-avatar.jpg"}`))
	request , _ := http.NewRequest("POST", "/auth/setavator", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+tokenString)

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, Error , response.Code)

	expectResponse := AvatorResponse{
		Code: Error,
		Msg:  "更新头像失败",
	}
	var actualResponse AvatorResponse
	json.Unmarshal(response.Body.Bytes(), &actualResponse)
	assert.Equal(t, expectResponse, actualResponse)
}

//正确的更换头像的用例

func TestSetAvator_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	//创建一个路由
	router := setupRouter()

	/* //生成一个token
	token, err := generateToken("testuser")
	assert.Equal(t,nil,err) */
	// 创建一个有效的Token
	claims := &middlewire.Claims{
		Username: "qin", // 根据自己数据库已有的用户
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(), // 设置有效的过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(middlewire.JwtKey)
	
	//创建一个POST请求
	body := bytes.NewBuffer([]byte(`{"url":"https://example.com/new-avatar.jpg"}`))
	request , _ := http.NewRequest("POST", "/auth/setavator", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+tokenString)

	//创建一个响应器
	response := httptest.NewRecorder()

	//执行请求
	router.ServeHTTP(response, request)

	//检查响应码
	assert.Equal(t, http.StatusOK, response.Code)

	//检查响应体
	expectResponse := AvatorResponse{
		Code: Success,
		Msg:  "获取头像成功",
		Data: "https://example.com/new-avatar.jpg",
	}
	var actualResponse AvatorResponse
	json.Unmarshal(response.Body.Bytes(), &actualResponse)
	assert.Equal(t, expectResponse, actualResponse)
}