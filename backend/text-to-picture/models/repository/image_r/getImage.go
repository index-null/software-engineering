package image_r

import (
	"net/http"
	"strconv"
	i "text-to-picture/models/image"
	d "text-to-picture/models/init"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 根据图片ID查询相关图片
func GetImagesById(db *gorm.DB, id int) (*i.ImageInformation, error) {
	var image i.ImageInformation
	err := db.Table("imageinformation").Where("id = ?", id).First(&image).Error
	if err != nil {
		return nil, err // 返回错误
	}

	return &image, nil // 返回指向image的指针
}

// 根据图片的username查询图片
func GetImagesByUsername(db *gorm.DB, username string) (*i.ImageInformation, error) {
	var image i.ImageInformation
	err := db.Table("imageinformation").Where("username = ?", username).First(&image).Error
	if err != nil {
		return nil, err // 返回错误
	}

	return &image, nil // 返回指向image的指针
}

func GetImages(c *gin.Context) {
	username := c.Query("username") // 从请求中获取用户名
	imageIdStr := c.Query("id")      // 从请求中获取图片ID（字符串）
	imageId, err := strconv.Atoi(imageIdStr) // 将字符串转换为整数

	if username != "" {
		image, err := GetImagesByUsername(d.DB, username)
		if err != nil {
			// 检查错误类型
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"message": "未找到相关图片"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "查询用户的图片失败", "err": "可能无该用户名的image"})
			}
			return
		}
		c.JSON(http.StatusOK, image)

	} else if err == nil {
		image, err := GetImagesById(d.DB, imageId)
		if err != nil {
			// 检查错误类型
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"message": "未找到相关图片"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "查询用户的图片失败", "err": "可能id不存在"})
			}
			return
		}
		c.JSON(http.StatusOK, image)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的用户ID或用户名"})
	}
}
