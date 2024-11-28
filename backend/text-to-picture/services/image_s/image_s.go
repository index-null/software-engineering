package image_s

import (
	"fmt"
	"text-to-picture/models"
	db "text-to-picture/models/init"
	"time"

	"gorm.io/gorm"
)

// GetImagesByUsername 获取用户的图片列表
func GetImagesByUsername(username string, offset, limit int) ([]models.Image, error) {
	var images []models.Image
	// 查询数据库，根据用户名获取图片列表
	err := db.Preload("User").Where("username = ?", username).Limit(limit).Offset(offset * limit).Find(&images).Error
	if err != nil {
		return nil, fmt.Errorf("查询图片失败: %v", err)
	}
	return images, nil
}

// LikeImage 处理用户点赞图片的逻辑
func LikeImage(imageID, username string) error {
	// 根据 username 查询用户ID
	var user models.User
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return fmt.Errorf("用户未找到: %v", err)
	}

	// 检查用户是否已经点赞过该图片
	var like models.ImageLike
	if err := db.Where("user_id = ? AND image_id = ?", user.ID, imageID).First(&like).Error; err == nil {
		return fmt.Errorf("已点赞") // 如果已经点赞，返回错误信息
	}

	// 创建新的点赞记录
	newLike := models.ImageLike{
		UserID:     user.ID,
		ImageID:    imageID,
		CreateTime: time.Now(),
	}
	if err := db.Create(&newLike).Error; err != nil {
		return err // 插入点赞记录失败，返回错误
	}

	// 更新图片的点赞数
	if err := db.Model(&models.Image{}).Where("id = ?", imageID).Update("likes_count", gorm.Expr("likes_count + ?", 1)).Error; err != nil {
		return err // 更新点赞数失败，返回错误
	}

	return nil // 点赞成功
}
