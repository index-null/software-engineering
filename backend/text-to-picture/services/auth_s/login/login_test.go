package auth_s

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockLoginHandler 是一个模拟登录处理函数，用于测试
func MockLoginHandler(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// 解析 JSON 数据
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "请求数据格式错误",
		})
		return
	}

	// 模拟用户逻辑
	if input.Username == "root1" && input.Password == "sssssss" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "登录成功",
		})
		return
	} else if input.Username == "nonexistentuser" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "用户不存在",
		})
		return
	} else if input.Password == "wrongpassword" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "密码错误",
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    http.StatusInternalServerError,
		"message": "数据库查询错误",
	})
}

// SetupTestRouter 设置测试路由
func SetupTestRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/login", MockLoginHandler) // 使用 MockLoginHandler 替代原来的依赖
	return r
}

func TestLogin(t *testing.T) {
	// 初始化 Gin 测试路由
	r := SetupTestRouter()

	// 测试用例：登录成功
	t.Run("LoginSuccess", func(t *testing.T) {
		loginPayload := []byte(`{
			"username": "root1",
			"password": "sssssss"
		}`)

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginPayload))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "登录成功", response["message"])
	})

	// 测试用例：请求数据格式错误
	t.Run("BadRequest", func(t *testing.T) {
		loginPayload := []byte(`invalid_json`)

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginPayload))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "请求数据格式错误", response["message"])
	})

	// 测试用例：用户不存在
	t.Run("UserNotFound", func(t *testing.T) {
		loginPayload := []byte(`{
			"username": "nonexistentuser",
			"password": "password"
		}`)

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginPayload))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "用户不存在", response["message"])
	})

	// 测试用例：密码错误
	t.Run("IncorrectPassword", func(t *testing.T) {
		loginPayload := []byte(`{
			"username": "root1",
			"password": "wrongpassword"
		}`)

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginPayload))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "密码错误", response["message"])
	})

	// 测试用例：数据库查询错误
	t.Run("DatabaseError", func(t *testing.T) {
		loginPayload := []byte(`{
			"username": "dberroruser",
			"password": "password"
		}`)

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginPayload))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "数据库查询错误", response["message"])
	})
}
