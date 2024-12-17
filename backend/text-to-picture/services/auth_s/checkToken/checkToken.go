package checktoken

import (
	"net/http"
	middlewire "text-to-picture/middlewire/jwt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

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

// 对token的情况进行响应
func CheckToken(c *gin.Context) {

	//获取请求头中的token
	tokenStr := c.GetHeader("Authorization")

	// 调用token校验函数
	isValid, msg := TokenCheck(tokenStr)

	// 返回错误情况的响应
	if !isValid {
		c.JSON(http.StatusUnauthorized, CheckTokenResponse{
			Code: http.StatusUnauthorized,
			Msg:  msg,
		})
		return
	}

	//返回正确情况的响应
	c.JSON(http.StatusOK, CheckTokenResponse{
		Code: http.StatusOK,
		Msg:  "令牌有效",
		Data: tokenStr,
	})
}

// 对token进行校验
func TokenCheck(tokenStr string) (bool, string) {

	// 解析token
	claims := &middlewire.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return middlewire.JwtKey, nil // jwtKey 是你的签名密钥
	})

	// 处理错误
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return false, "令牌格式不正确"
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return false, "令牌过期或未激活"
			} else {
				return false, "令牌无法处理"
			}
		}
		return false, "令牌无法处理"
	}

	if !token.Valid {
		return false, "令牌无效"
	}

	return true, "令牌有效"
}
