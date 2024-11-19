package image_r

import (
	"net/http"
	"strconv"
	i "text-to-picture/models/image"
	d "text-to-picture/models/init"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 根据用户ID查询相关图片
func GetImagesByUserId(db *gorm.DB, userId int) ([]i.ImageInformation, error) {
	var images []i.ImageInformation
	err := db.Table("imageinformation").Where("id = ?", userId).Find(&images).Error // 使用 Find 而不是 First
	if err != nil {
		return nil, err
	}

	return images, nil
}

// 根据用户ID查询收藏的图片
func GetFavoritedImagesByUserId(db *gorm.DB, userId int) ([]i.ImageInformation, error) {
	var images []i.ImageInformation
	err := db.Table("favoritedimage").Where("id = ?", userId).Find(&images).Error // 使用 Find 而不是 First
	if err != nil {
		return nil, err
	}

	return images, nil
}

// 根据用户名查询相关图片
func GetImagesByUsername(db *gorm.DB, username string) ([]i.ImageInformation, error) {
	var images []i.ImageInformation
	err := db.Table("imageinformation").Where("username = ?", username).Find(&images).Error // 使用 Find 而不是 First
	if err != nil {
		return nil, err
	}

	return images, nil
}

// 根据用户名查询收藏的图片
func GetFavoritedImagesByUsername(db *gorm.DB, username string) ([]i.ImageInformation, error) {
	var images []i.ImageInformation
	err := db.Table("favoritedimage").Where("username = ?", username).Find(&images).Error // 使用 Find 而不是 First
	if err != nil {
		return nil, err
	}

	return images, nil
}
//-------------------------------------------------------------------------------------------
func GetUserImagesById(c *gin.Context) {
	userIdStr := c.Query("id") // 从请求中获取用户ID（字符串）
	userId, err := strconv.Atoi(userIdStr) // 将字符串转换为整数
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的用户ID"})
		return
	}

	images, err := GetImagesByUserId(d.DB, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询用户图片失败"})
		return
	}

	c.JSON(http.StatusOK, images)
}

func GetUserFavoritedImagesById(c *gin.Context) {
	userIdStr := c.Query("id") // 从请求中获取用户ID（字符串）
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

func GetUserImagesByUsername(c *gin.Context) {
	username := c.Query("username") // 从请求中获取用户名
	//c.JSON(200,gin.H{"name":username})
	images, err := GetImagesByUsername(d.DB, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询用户图片失败"})
		return
	}

	c.JSON(http.StatusOK, images)
}

func GetUserFavoritedImagesByUsername(c *gin.Context) {
	username := c.Query("username") // 从请求中获取用户名

	images, err := GetFavoritedImagesByUsername(d.DB, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询用户收藏的图片失败"})
		return
	}

	c.JSON(http.StatusOK, images)
}

