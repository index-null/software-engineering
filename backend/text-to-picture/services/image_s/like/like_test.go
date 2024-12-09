package like

import (
	"bytes"
	"encoding/json"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"text-to-picture/models/image"
	db "text-to-picture/models/init"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLikeImage(t *testing.T) {
	// 初始化 Gin 引擎
	r := gin.Default()
	r.PUT("/auth/like", LikeImage)

	// 模拟数据库
	db.DB = MockDB()

	// 测试用例
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		username       string
		mockDBSetup    func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "成功点赞",
			requestBody:    map[string]interface{}{"url": "http://example.com/image.jpg"},
			username:       "test_user",
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"current_likes": 6, "message": "Image liked successfully"},
		},
		{
			name:           "缺少图片 URL",
			requestBody:    map[string]interface{}{},
			username:       "test_user",
			mockDBSetup:    func() {}, // 不需要数据库操作
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"code": 400, "error": "Missing image URL"},
		},
		{
			name:        "用户已点赞",
			requestBody: map[string]interface{}{"url": "http://example.com/image.jpg"},
			username:    "test_user",
			mockDBSetup: func() {
				// 初始化用户点赞记录
				db.DB.Create(&image.ImageLike{Picture: "http://example.com/image.jpg", UserName: "test_user", Num: 1})
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   map[string]interface{}{"code": 409, "error": "用户已经点赞过该图片"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// 初始化数据库状态
			test.mockDBSetup()

			// 模拟请求
			body, _ := json.Marshal(test.requestBody)
			req := httptest.NewRequest("PUT", "/auth/like", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer token") // 模拟 Token

			// 模拟上下文
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("username", test.username) // 模拟上下文中解析到的用户名

			// 执行路由
			r.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, test.expectedStatus, w.Code)

			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.Equal(t, test.expectedBody, responseBody)
		})
	}
}

// MockDB 模拟数据库操作
func MockDB() *gorm.DB {
	// 使用内存数据库进行测试
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&image.ImageLike{}, &image.ImageInformation{})
	return db
}
