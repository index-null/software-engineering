package avator

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"text-to-picture/middlewire/jwt"
	"text-to-picture/models/init"
	"text-to-picture/models/user"
)

type AvatorResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

const (
	Success      = 200
	Error        = 500
	Unauthorized = 401
)

func SetAvator(c *gin.Context) {
	tokenStr := c.GetHeader("Authorization")
	newURL := c.Query("url")

	if tokenStr == "" {
		c.JSON(Unauthorized, AvatorResponse{
			Code: Unauthorized,
			Msg:  "请求头中缺少Token",
		})
		return
	}

	claims := &middlewire.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return middlewire.JwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(Unauthorized, AvatorResponse{
			Code: Unauthorized,
			Msg:  "无效的Token",
		})
		return
	}

	username := claims.Username

	// 更新数据库中的头像 URL
	result := models.DB.Model(&user.UserInformation{}).Where("username = ?", username).Update("avatar_url", newURL)
	if result.Error != nil {
		c.JSON(Error, AvatorResponse{
			Code: Error,
			Msg:  "更新头像失败",
		})
		return
	}

	c.JSON(Success, AvatorResponse{
		Code: Success,
		Msg:  "头像更新成功",
		Data: newURL,
	})
}
func GetAvator(c *gin.Context) {
	tokenStr := c.GetHeader("Authorization")

	if tokenStr == "" {
		c.JSON(Unauthorized, AvatorResponse{
			Code: Unauthorized,
			Msg:  "请求头中缺少Token",
		})
		return
	}

	claims := &middlewire.Claims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return middlewire.JwtKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired > 0 {
				c.JSON(Unauthorized, AvatorResponse{
					Code: Unauthorized,
					Msg:  "Token已过期",
				})
				return
			}
		}
		c.JSON(Unauthorized, AvatorResponse{
			Code: Unauthorized,
			Msg:  "无效的Token",
		})
		return
	}

	username := claims.Username

	// 查询数据库中的头像 URL
	var usera user.UserInformation
	result := models.DB.Where("username = ?", username).First(&usera)
	if result.Error != nil {
		c.JSON(Error, AvatorResponse{
			Code: Error,
			Msg:  "查询头像失败",
		})
		return
	}

	c.JSON(Success, AvatorResponse{
		Code: Success,
		Msg:  "获取头像成功",
		Data: usera.Avatar_url,
	})
}
