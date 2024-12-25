package checktoken

import (
	"net/http"
	middlewire "text-to-picture/middlewire/jwt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// CheckTokenResponse 是一个结构体，用于定义检查token接口的响应格式
type CheckTokenResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// @Summary 校验token是否有效
// @Description 校验token接口
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} CheckTokenResponse "令牌有效"
// @Failure 401 {object} CheckTokenResponse "令牌格式错误"
// @Failure 401 {object} CheckTokenResponse "令牌过期或者未激活"
// @Failure 401 {object} CheckTokenResponse "令牌无法处理"
// @Failure 401 {object} CheckTokenResponse "令牌无效"
// @Router /auth/checkToken [get]
// CheckToken 函数用于校验请求中的token是否有效
func CheckToken(c *gin.Context) {
	// 获取请求头中的token
	tokenStr := c.GetHeader("Authorization")
	// 调用token校验函数
	username, isValid, msg := TokenCheck(tokenStr)
	// 返回错误情况的响应
	if !isValid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": msg,
		})
		return
	}
	// 检查用户名是否为root
	if username != "root" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "不是root用户",
		})
		return
	}
	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "令牌有效,识别为至高无上的root用户",
	})
}

// TokenCheck 函数用于校验token的有效性
// 它接受一个token字符串作为输入，返回用户名、token是否有效以及错误信息
func TokenCheck(tokenStr string) (string, bool, string) {
	// 解析token
	claims := &middlewire.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return middlewire.JwtKey, nil // jwtKey 是你的签名密钥
	})
	username := claims.Username
	// 处理错误
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			// token格式错误
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return username, false, "令牌格式不正确"
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// token过期或未激活
				return username, false, "令牌过期或未激活"
			} else {
				// 其他token处理错误
				return username, false, "令牌无法处理"
			}
		}
		// 其他错误
		return username, false, "令牌无法处理"
	}
	// token无效
	if !token.Valid {
		return "", false, "令牌无效"
	}
	// token有效
	return username, true, "令牌有效"
}
