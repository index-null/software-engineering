package middlewire

import (
	"fmt"
	"log"
	models "text-to-picture/models/init"
	"text-to-picture/models/user"
	u "text-to-picture/models/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func AddScore(c *gin.Context) {
	// 从上下文中获取用户名
	Username, exists := c.Get("username")
	if !exists {
		log.Printf("未找到用户名")
		c.JSON(401, gin.H{
			"code":    401,
			"success": false,
			"message": "未找到用户信息",
		})
		return
	}
	username := Username.(string)
	// 更新用户积分
	result := models.DB.Model(&user.UserInformation{}).Where("username = ?", username).Update("score", gorm.Expr("score + ?", 100))
	if result.Error != nil {
		c.JSON(401, gin.H{
			"code":    401,
			"success": false,
			"message": "积分记录更新失败",
		})
		return
	}

	var user u.UserInformation
	if err := models.DB.Where("username = ?", username).First(&user).Error; err != nil {
		log.Printf("Failed to query user information: %v", err)
		c.JSON(401, gin.H{
			"code":    401,
			"success": false,
			"message": "用户信息查询失败",
		})
	}
	var record u.UserScore
	record.Username = username
	record.Record = "积分+100"
	if err := models.DB.Create(&record).Error; err != nil {
		log.Printf("Failed to create record: %v", err)
		c.JSON(401, gin.H{
			"code":    401,
			"success": false,
			"message": "积分记录创建失败",
		})
	}
	msg := fmt.Sprintf("积分记录更新成功!用户当前积分为%v", user.Score)
	c.JSON(200, gin.H{
		"code":    200,
		"success": true,
		"message": msg,
	})
}