package image_r

import (
	"net/http"
	"strconv"
	i "text-to-picture/models/image"
	d "text-to-picture/models/init"
	u "text-to-picture/models/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 根据用户ID查询相关图片
func GetUserImagesByUserId(db *gorm.DB, userId int) ([]i.ImageInformation, error) {
	var images []i.ImageInformation
	err := db.Table("imageinformation").Where("user.id = ?", userId).Find(&images).Error // 使用 Find 而不是 First
	if err != nil {
		return nil, err
	}

	return images, nil
}

// 根据用户ID查询收藏的图片
func GetUserFavoritedImagesByUserId(db *gorm.DB, userId int) ([]i.ImageInformation, error) {
	var images []i.ImageInformation
	err := db.Table("favoritedimage").Where("user.id = ?", userId).Find(&images).Error // 使用 Find 而不是 First
	if err != nil {
		return nil, err
	}

	return images, nil
}

// 根据用户名查询相关图片
func GetUserImagesByUsername(db *gorm.DB, username string) ([]i.ImageInformation, error) {
	var images []i.ImageInformation
	err := db.Table("imageinformation").Where("username = ?", username).Find(&images).Error // 使用 Find 而不是 First
	if err != nil {
		return nil, err
	}

	return images, nil
}

// 根据用户名查询收藏的图片
func GetUserFavoritedImagesByUsername(db *gorm.DB, username string) ([]i.ImageInformation, error) {
	var images []i.ImageInformation
	err := db.Table("favoritedimage").Where("username = ?", username).Find(&images).Error // 使用 Find 而不是 First
	if err != nil {
		return nil, err
	}

	return images, nil
}

//-------------------------------------------------------------------------------------------
//目前的imageinformation表并没有直接的userId属性

func GetUserImages(c *gin.Context) {
	username := c.Query("username") // 从请求中获取用户名
	userIdStr := c.Query("id") // 从请求中获取用户ID（字符串）
	userId, err := strconv.Atoi(userIdStr) // 将字符串转换为整数
	
	if username != ""{
		images, err := GetUserImagesByUsername(d.DB, username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "查询用户图片失败","err":err})
			return
		}
		c.JSON(http.StatusOK, images)
		return

	}else if err == nil{// id转username
		var user u.UserInformation
		err := d.DB.Table("userinformation").Where("id = ?", userId).First(&user).Error // 使用 Find 而不是 First
			if err != nil {
				images, err := GetUserImagesByUserId(d.DB, user.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": "查询用户图片失败","err":err})
					return
				}
				c.JSON(http.StatusOK, images)
				return
			}else{
				c.JSON(http.StatusBadRequest, gin.H{"message": "无效的用户id"})
			}

	}else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的用户名或用户id"})
		return
	}
}

func GetUserFavoritedImages(c *gin.Context) {
	username := c.Query("username") // 从请求中获取用户名
	userIdStr := c.Query("id") // 从请求中获取用户ID（字符串）
	userId, err := strconv.Atoi(userIdStr) // 将字符串转换为整数

	if username != ""{
		images, err := GetUserFavoritedImagesByUsername(d.DB, username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "查询用户收藏的图片失败","err":err})
			return
		}
		c.JSON(http.StatusOK, images)
		return

	}else if err == nil{// id转username
		var user u.UserInformation
		err := d.DB.Table("userinformation").Where("id = ?", userId).First(&user).Error // 使用 Find 而不是 First
			if err != nil {
				images, err := GetUserFavoritedImagesByUserId(d.DB, user.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": "查询用户收藏的图片失败","err":err})
					return
				}
				c.JSON(http.StatusOK, images)
				return
			}else{
				c.JSON(http.StatusBadRequest, gin.H{"message": "无效的用户id"})
			}

	}else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的用户ID或用户名"})
		return
	}
	
}

	// else if err == nil{
	// 	images, err := GetUserImagesByUserId(d.DB, userId)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询用户图片失败","err":err})
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, images)
	// 	return
	// }

// else if err == nil{
	// 	images, err := GetUserFavoritedImagesByUserId(d.DB, userId)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询用户收藏的图片失败","err":err})
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, images)
	// 	return

	// }
