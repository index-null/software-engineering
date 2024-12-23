package update

import (
	"fmt"
	"log"
	models "text-to-picture/models/init"
	u "text-to-picture/models/user"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddScore(c *gin.Context) {
	fmt.Println("_------------------------------------")
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

	// 获取当前时间
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	todayEnd := todayStart.Add(24 * time.Hour)

	var existingRecord u.UserScore
	// 检查今天是否已有签到记录
	result := models.DB.Where("username = ? AND create_time >= ? AND create_time < ?", username, todayStart, todayEnd).First(&existingRecord)
	if result.Error == nil {
		// 如果找到了今天的记录，则认为用户已经签到过
		c.JSON(400, gin.H{
			"code":    400,
			"success": false,
			"message": "您今天已经签到过了",
		})
		return
	} else if result.Error != gorm.ErrRecordNotFound {
		// 如果查询时发生其他错误
		log.Printf("Failed to check sign-in record: %v", result.Error)
		c.JSON(500, gin.H{
			"code":    500,
			"success": false,
			"message": "检查签到记录失败",
		})
		return
	}

	// 更新用户积分
	updateResult := models.DB.Model(&u.UserInformation{}).Where("username = ?", username).Update("score", gorm.Expr("score + ?", 100))
	if updateResult.Error != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"success": false,
			"message": "积分记录更新失败",
		})
		return
	}

	// 创建新的签到记录
	newRecord := u.UserScore{
		Username:   username,
		Record:     "积分+100",
		Create_time: time.Now(),
	}
	if err := models.DB.Create(&newRecord).Error; err != nil {
		log.Printf("Failed to create sign-in record: %v", err)
		c.JSON(500, gin.H{
			"code":    500,
			"success": false,
			"message": "创建签到记录失败",
		})
		return
	}

	// 查询更新后的用户信息以获取最新积分
	var updatedUser u.UserInformation
	if err := models.DB.Where("username = ?", username).First(&updatedUser).Error; err != nil {
		log.Printf("Failed to query updated user information: %v", err)
		c.JSON(500, gin.H{
			"code":    500,
			"success": false,
			"message": "用户信息查询失败",
		})
		return
	}

	msg := fmt.Sprintf("积分记录更新成功!用户当前积分为%v", updatedUser.Score)
	c.JSON(200, gin.H{
		"code":    200,
		"success": true,
		"message": msg,
	})
}