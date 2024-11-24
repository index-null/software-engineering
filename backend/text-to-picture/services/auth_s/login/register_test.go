// register_test.go
package auth_s

import (
	"bytes"
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"text-to-picture/models/user"
)

// MockUserRepository 是一个模拟的用户仓库
type MockUserRepository struct{}

// InsertUserInformation 是一个模拟的插入用户信息方法
func (m *MockUserRepository) InsertUserInformation(user *user.UserInformation) error {
	// 模拟插入成功的场景
	return nil
}

// GetUserByName 是一个模拟的获取用户方法
func (m *MockUserRepository) GetUserByName(name string) (*user.UserInformation, error) {
	// 模拟用户不存在的场景
	return nil, gorm.ErrRecordNotFound
}

// SetupRouter 设置 Gin 路由
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/register", Register)
	return r
}

// TestRegister_Success 测试注册成功的情况
func TestRegister_Success(t *testing.T) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 创建一个新的路由器
	router := SetupRouter()

	// 创建一个 POST 请求
	body := bytes.NewBuffer([]byte(`{"email": "test@example.com", "username": "testuser", "password": "testpassword"}`))
	req, _ := http.NewRequest("POST", "/register", body)
	req.Header.Set("Content-Type", "application/json")

	// 创建一个响应记录器
	w := httptest.NewRecorder()

	// 执行请求
	router.ServeHTTP(w, req)

	// 检查响应状态码
	assert.Equal(t, http.StatusOK, w.Code)

	// 检查响应体
	expectedResponse := gin.H{
		"code":    http.StatusOK,
		"message": "注册成功",
	}
	var actualResponse gin.H
	json.Unmarshal(w.Body.Bytes(), &actualResponse)
	assert.Equal(t, expectedResponse, actualResponse)
}

// TestRegister_BadRequest 测试请求数据格式错误的情况
func TestRegister_BadRequest(t *testing.T) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 创建一个新的路由器
	router := SetupRouter()

	// 创建一个 POST 请求，故意发送格式错误的数据
	body := bytes.NewBuffer([]byte(`{"email": "test@example.com", "username": "testuser"}`)) // 缺少 password 字段
	req, _ := http.NewRequest("POST", "/register", body)
	req.Header.Set("Content-Type", "application/json")

	// 创建一个响应记录器
	w := httptest.NewRecorder()

	// 执行请求
	router.ServeHTTP(w, req)

	// 检查响应状态码
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 检查响应体
	expectedResponse := gin.H{
		"code":    http.StatusBadRequest,
		"message": "请求数据格式错误",
	}
	var actualResponse gin.H
	json.Unmarshal(w.Body.Bytes(), &actualResponse)
	assert.Equal(t, expectedResponse, actualResponse)
}

// TestRegister_InternalServerError 测试用户创建失败的情况
func TestRegister_InternalServerError(t *testing.T) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 创建一个新的路由器
	router := SetupRouter()

	// 创建一个 POST 请求
	body := bytes.NewBuffer([]byte(`{"email": "test@example.com", "username": "testuser", "password": "testpassword"}`))
	req, _ := http.NewRequest("POST", "/register", body)
	req.Header.Set("Content-Type", "application/json")

	// 创建一个响应记录器
	w := httptest.NewRecorder()

	// 执行请求
	router.ServeHTTP(w, req)

	// 检查响应状态码
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 检查响应体
	expectedResponse := gin.H{
		"code":    http.StatusInternalServerError,
		"message": "用户创建失败",
		"error":   "数据库错误",
	}
	var actualResponse gin.H
	json.Unmarshal(w.Body.Bytes(), &actualResponse)
	assert.Equal(t, expectedResponse, actualResponse)
}
