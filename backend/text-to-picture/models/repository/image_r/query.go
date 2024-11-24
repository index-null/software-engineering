package image_r

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	i "text-to-picture/models/image"
	d "text-to-picture/models/init"
	u "text-to-picture/models/user"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// // 根据用户ID查询相关图片
// func GetUserImagesByUserId(db *gorm.DB, userId int) ([]i.ImageInformation, error) {
// 	var images []i.ImageInformation
// 	err := db.Table("imageinformation").Where("userid = ?", userId).Find(&images).Error // 使用 Find 而不是 First
// 	if err != nil {
// 		return nil, err
// 	}

// 	return images, nil
// }

// // 根据用户ID查询收藏的图片
// func GetUserFavoritedImagesByUserId(db *gorm.DB, userId int) ([]i.ImageInformation, error) {
// 	var images []i.ImageInformation
// 	err := db.Table("favoritedimage").Where("userid = ?", userId).Find(&images).Error // 使用 Find 而不是 First
// 	if err != nil {
// 		return nil, err
// 	}

// 	return images, nil
// }

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

//查询指定时间段内的所有图像
func GetImagesInfoWithinTimeRange(db *gorm.DB, startTime, endTime time.Time) ([]i.ImageInformation, error) {
	var images []i.ImageInformation
	err := db.Table("imageinformation").
		Where("create_time BETWEEN ? AND ?", startTime, endTime).
		Find(&images).Error
	if err != nil {
		return nil, fmt.Errorf("查询图像列表时发生错误: %v", err)
	}
	return images, nil
}

// 获取所有图像信息并按id排序
func GetAllImagesInfo(db *gorm.DB) ([]i.ImageInformation, error) {
	var images []i.ImageInformation
	result := db.Table("imageinformation").Order("id ASC").Find(&images)
	if result.Error != nil {
		return nil, fmt.Errorf("查询图像列表时发生错误: %v", result.Error)
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
				c.JSON(http.StatusBadRequest, gin.H{"message": "无效的用户id"})
			}
			
			images, err := GetUserImagesByUsername(d.DB, user.UserName)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "查询用户图片失败","err":err})
				return
			}
			c.JSON(http.StatusOK, images)
			return

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
				c.JSON(http.StatusBadRequest, gin.H{"message": "无效的用户id"})
			}

			images, err := GetUserFavoritedImagesByUsername(d.DB, user.UserName)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "查询用户收藏的图片失败","err":err})
				return
			}
			c.JSON(http.StatusOK, images)
			return

	}else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的用户ID或用户名"})
		return
	}
	
}

//查询指定时间段内的所有图像
func GetImagesWithinTimeRange(c *gin.Context) {
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	//定义正则表达式来检测时间字符串是否包含时间部分
	timeRegex := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2})(\.\d+)?Z$`)
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

	var startTime, endTime time.Time
	var err error

	//检查 start_time 是否包含时间部分
	if timeRegex.MatchString(startTimeStr) {//含时间部分
		startTime, err = time.Parse("2006-01-02T15:04:05Z", startTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "无效的开始时间格式", "error": err.Error()})
			return
		}
	} else if dateRegex.MatchString(startTimeStr) {//不含时间部分
		startTime, err = time.Parse("2006-01-02", startTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "无效的开始时间格式", "error": err.Error()})
			return
		}
		startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, time.UTC)
	} else {//都不符合
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的开始时间格式"})
		return
	}

	//检查 end_time 是否包含时间部分
	if timeRegex.MatchString(endTimeStr) {//含时间部分
		endTime, err = time.Parse("2006-01-02T15:04:05Z", endTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "无效的结束时间格式", "error": err.Error()})
			return
		}
	} else if dateRegex.MatchString(endTimeStr) {//不含时间部分
		endTime, err = time.Parse("2006-01-02", endTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "无效的结束时间格式", "error": err.Error()})
			return
		}
		endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 999999999, time.UTC)
	} else {//都不符合
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的结束时间格式"})
		return
	}

	images, err := GetImagesInfoWithinTimeRange(d.DB, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询图像列表失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "查询图像列表成功",
		"images":  images,
	})
}

// 获取所有图像信息
func GetAllImages(c *gin.Context) {
	images, err := GetAllImagesInfo(d.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "获取图像列表失败", "error": err.Error()})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"message": "获取图像列表成功",
		"images":  images,
	})
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
