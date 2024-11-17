package image_u

import (
	"net/http"
	"strconv"
	i "text-to-picture/models/image"
	d "text-to-picture/models/init"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 根据用户ID查询相关图片
func GetImagesByUserId(db *gorm.DB, userId int) ([]i.Image, error) {
	var images []i.Image
	err := db.Table("images").Where("user_id = ?", userId).Find(&images).Error // 使用 Find 而不是 First
	if err != nil {
		return nil, err
	}

	return images, nil
}

// 根据用户ID查询收藏的图片
func GetFavoritedImagesByUserId(db *gorm.DB, userId int) ([]i.Image, error) {
	var images []i.Image
	err := db.Table("FavoritedImage").Where("user_id = ?", userId).Find(&images).Error // 使用 Find 而不是 First
	if err != nil {
		return nil, err
	}

	return images, nil
}

func GetUserImages(c *gin.Context) {
	userIdStr := c.Param("user_id") // 从请求中获取用户ID（字符串）
	userId, err := strconv.Atoi(userIdStr) // 将字符串转换为整数

	images, err := GetImagesByUserId(d.DB, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询用户图片失败"})
		return
	}

	c.JSON(http.StatusOK, images)
}

func GetUserFavoritedImages(c *gin.Context) {
	userIdStr := c.Param("user_id") // 从请求中获取用户ID（字符串）
	userId, err := strconv.Atoi(userIdStr) // 将字符串转换为整数
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的用户ID"})
		return
	}

	images, err := GetFavoritedImagesByUserId(d.DB, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询用户收藏的图片失败"})
		return
	}

	c.JSON(http.StatusOK, images)
}