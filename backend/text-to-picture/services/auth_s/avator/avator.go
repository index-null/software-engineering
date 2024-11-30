package avator

import (
	"github.com/gin-gonic/gin"
	"net/http"
	models "text-to-picture/models/init"
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
	var reqBody struct {
		URL string `json:"url"`
	}

	c.BindJSON(&reqBody)

	newURL := reqBody.URL
	usernames, _ := c.Get("username")
	if usernames == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":  Unauthorized,
			"error": "名字解析出错"})
		return
	}
	username, _ := usernames.(string)

	// 更新数据库中的头像 URL
	result := models.DB.Model(&user.UserInformation{}).Where("username = ?", username).Update("avatar_url", newURL)
	if result.Error != nil || result.RowsAffected == 0 {
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
	username, _ := c.Get("username")
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":  Unauthorized,
			"error": "名字解析出错"})
		return
	}

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
