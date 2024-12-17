package delete

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	middlewire "text-to-picture/middlewire/jwt"
	db "text-to-picture/models/init"
	u "text-to-picture/models/user"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type DeleteUserTestCase struct {
	Name            string
	Username        string
	TokenUsername   string
	IsOwn           string
	ExpectedCode    int
	ExpectedMessage string
	needCreate		bool
}

var testCases = []DeleteUserTestCase{
	{
		Name:            "TestDeleteUserByName_Success",
		Username:        "testuser_success3",
		TokenUsername:   "root",
		ExpectedCode:    http.StatusOK,
		ExpectedMessage: "成功删除用户：testuser_success3",
		needCreate:		 true,
	},
	{
		Name:            "TestDeleteUserByName_LogoutSuccess",
		Username:        "testuser_logout5",
		TokenUsername:   "testuser_logout5",
		IsOwn:           "true",
		ExpectedCode:    http.StatusOK,
		ExpectedMessage: "testuser_logout5的账号注销成功",
		needCreate:		 true,
	},
	{
		Name:            "TestDeleteUserByName_UserNotFound",
		Username:        "nonexistentuser1",
		TokenUsername:   "root",
		ExpectedCode:    http.StatusNotFound,
		ExpectedMessage: "用户不存在",
		needCreate:		 false,
	},
	{
		Name:            "TestDeleteUserByName_NoPermission",
		Username:        "testuser_non2",
		TokenUsername:   "nonrootuser",
		ExpectedCode:    http.StatusBadRequest,
		ExpectedMessage: "非root用户，不可删除其他某个用户",
		needCreate:		 true,
	},
	{
		Name:            "TestDeleteUserByName_NoToken",
		Username:        "notoken",
		ExpectedCode:    http.StatusUnauthorized,
		ExpectedMessage: "请求头中缺少Token",
		needCreate:		 false,
	},
	{
		Name:            "TestDeleteUserByName_InvalidToken",
		Username:        "testuser",
		TokenUsername:   "invalid-token",
		ExpectedCode:    http.StatusUnauthorized,
		ExpectedMessage: "无效的Token",
		needCreate:		 false,
	},
}

func TestDeleteUserByName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := SetupRouter()

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// 清理数据库，确保每次测试开始前状态一致
			cleanupDB()

			if tc.TokenUsername != "" {
				if tc.needCreate {
					// 创建测试用户
					db.DB.Create(&u.UserInformation{UserName: tc.Username, Email: tc.Username + "@qq.com"})
				}

				// 创建有效的Token或无效Token
				var tokenString string
				if tc.TokenUsername == "invalid-token" {
					tokenString = "invalid-token"
				} else {
					claims := &middlewire.Claims{
						Username: tc.TokenUsername,
						StandardClaims: jwt.StandardClaims{
							ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
						},
					}
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
					tokenString, _ = token.SignedString(middlewire.JwtKey)
				}

				// 构建请求URL
				url := fmt.Sprintf("/auth/root/deleteOneUser?username=%s", tc.Username)
				if tc.IsOwn != "" {
					url += fmt.Sprintf("&isOwn=%s", tc.IsOwn)
				}

				req, _ := http.NewRequest("DELETE", url, nil)
				req.Header.Set("Authorization", tokenString)

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				// 检查响应状态码
				assert.Equal(t, tc.ExpectedCode, w.Code)

				// 检查响应体
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				if response["message"] != nil {
					assert.Equal(t, tc.ExpectedMessage, response["message"].(string))
				} else {
					// 处理无消息返回的情况，例如401未授权时可能没有具体的消息
					assert.Equal(t, tc.ExpectedMessage, "")
				}
			} else {
				// 测试无Token情况
				req, _ := http.NewRequest("DELETE", "/auth/root/deleteOneUser?username="+tc.Username, nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				// 检查响应状态码
				assert.Equal(t, tc.ExpectedCode, w.Code)

				// 检查响应体
				var response gin.H
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Equal(t, tc.ExpectedMessage, response["message"])
			}
		})
	}
}

// cleanupDB 函数用于清理数据库，在每次测试前调用以保证数据库状态的一致性
func cleanupDB() {
	// 实现清理逻辑，如删除所有测试用户等
	// 注意：在实际使用中应更加谨慎处理，以免误删生产数据
	db.DB.Where("username LIKE ?", "testuser_%").Delete(&u.UserInformation{})
}