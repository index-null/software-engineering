package query

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	d "text-to-picture/models/init"
	"text-to-picture/models/repository/image_r"

	//u "text-to-picture/models/user"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/*
	获取用户生成的图像			GetUserImages
	获取用户的收藏图像			GetUserFavoritedImages

	查询指定的某张图像			GetImage
	查询指定时间段内的所有图像	 GetImagesWithinTimeRange
	获取所有图像信息			GetAllImages
*/

// 获取当前登录用户生成的图像
// @Summary 获取当前登录用户生成的图像
// @Description 获取当前用户生成的所有图像
// @Tags images
// @Produce json
// @Success 200 {object} map[string]interface{} "获取用户的图像成功"
// @Failure 401 {object} map[string]interface{} "未找到用户信息"
// @Failure 500 {object} map[string]interface{} "查询用户图片失败"
// @Router /auth/user/images [get]
// GetUserImages 根据用户名获取用户的图片列表
// 参数: c *gin.Context，包含请求上下文和路由信息
func GetUserImages(c *gin.Context) {
	// 从上下文中获取用户名
	username, exists := c.Get("username")
	// 如果用户名不存在，则返回401错误响应
	if !exists {
		log.Printf("未找到用户名")
		c.JSON(401, gin.H{
			"success": false,
			"message": "未找到用户信息",
		})
		return
	}

	// 根据用户名查询用户图片
	images, err := image_r.GetUserImagesByUsername(d.DB, username.(string))
	// 如果查询失败，则返回500错误响应
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询用户图片失败", "error": err})
		return
	}
	// 返回200成功响应，包含用户图片列表
	c.JSON(http.StatusOK, gin.H{"message": "获取用户的图像成功", "images": images})
	return
}

// 获取当前登录用户的收藏图像
// @Summary 获取当前登录用户的收藏图像
// @Description 获取当前用户收藏的所有图像
// @Tags favorites
// @Produce json
// @Success 200 {object} map[string]interface{} "获取用户收藏的图像成功"
// @Failure 401 {object} map[string]interface{} "未找到用户信息"
// @Failure 500 {object} map[string]interface{} "查询用户收藏的图片失败"
// @Router /auth/user/favoritedimages [get]
// GetUserFavoritedImages 获取用户收藏的图片
// 该函数从上下文中提取用户名，然后查询该用户收藏的图片并返回
func GetUserFavoritedImages(c *gin.Context) {
	// 从上下文中获取用户名
	username, exists := c.Get("username")
	fmt.Println(username.(string))
	// 如果未找到用户名，则记录错误并返回错误响应
	if !exists {
		log.Printf("未找到用户名")
		c.JSON(401, gin.H{
			"success": false,
			"message": "未找到用户信息",
		})
		return
	}

	// 调用服务层函数，根据用户名查询用户收藏的图片
	images, err := image_r.GetUserFavoritedImagesByUsername(d.DB, username.(string))
	// 如果查询失败，则记录错误并返回错误响应
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询用户收藏的图片失败", "err": err})
		return
	}
	// 查询成功，返回用户收藏的图片列表
	c.JSON(http.StatusOK, images)
}

// 查询指定的某张图像
// @Summary 查询指定图像
// @Description 根据图像的 URL、用户名或 ID 查询某张图像
// @Tags images
// @Produce json
// @Param url query string false "图像的URL"
// @Param username query string false "用户名"
// @Param id query int false "图像ID"
// @Success 200 {object} map[string]interface{} "查询图像成功"
// @Failure 400 {object} map[string]interface{} "无效的图像ID或用户名"
// @Failure 404 {object} map[string]interface{} "未找到相关图片"
// @Failure 500 {object} map[string]interface{} "查询用户的图片失败"
// @Router /image [get]

// GetImage 根据请求参数获取图片信息。
// 该函数首先从请求参数中提取url、username和id，然后根据这些参数查询数据库以获取图片信息。
// 如果url、username或id为空或无效，函数将返回相应的错误信息。
// 如果查询成功，函数将返回图片信息。
func GetImage(c *gin.Context) {
	url := c.Query("url")
	username := c.Query("username")          // 从请求中获取用户名
	imageIdStr := c.Query("id")              // 从请求中获取图片ID（字符串）
	imageId, err := strconv.Atoi(imageIdStr) // 将字符串转换为整数

	if url != "" {
		image, err := image_r.GetImageByUrl(d.DB, url)
		if err != nil {
			// 检查错误类型
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"message": "未找到相关图片"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "查询用户的图片失败", "error": err})
			}
			return
		}
		c.JSON(http.StatusOK, image)
	} else if username != "" {
		image, err := image_r.GetImageByUsername(d.DB, username)
		if err != nil {
			// 检查错误类型
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"message": "未找到相关图片"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "查询用户的图片失败", "error": err})
			}
			return
		}
		c.JSON(http.StatusOK, image)

	} else if err == nil {
		image, err := image_r.GetImageById(d.DB, imageId)
		if err != nil {
			// 检查错误类型
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"message": "未找到相关图片"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "查询用户的图片失败", "error": err})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "查询图像成功", "image": image})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的图像ID或用户名"})
	}
}

// 查询当前登录用户在指定时间段内生成过的所有图像
// @Summary 查询用户在指定时间段内生成的所有图像
// @Description 获取当前用户在指定时间范围内生成的图像列表
// @Tags images
// @Produce json
// @Param start_time query string true "开始时间 (格式: YYYY-MM-DD 或 YYYY-MM-DDTHH:MM:SSZ)"
// @Param end_time query string true "结束时间 (格式: YYYY-MM-DD 或 YYYY-MM-DDTHH:MM:SSZ)"
// @Success 200 {object} map[string]interface{} "查询图像列表成功"
// @Failure 401 {object} map[string]interface{} "未找到用户信息"
// @Failure 400 {object} map[string]interface{} "无效的开始或结束时间格式"
// @Failure 500 {object} map[string]interface{} "查询图像列表失败"
// @Router /auth/user/images/timeRange [get]
// GetImagesWithinTimeRange 根据时间范围获取图像列表
// 该函数从gin上下文中获取用户信息，并根据查询参数中的开始和结束时间
// 从数据库中获取该时间范围内的图像信息列表
func GetImagesWithinTimeRange(c *gin.Context) {
	// 从上下文中获取用户名
	username, exists := c.Get("username")
	if !exists {
		log.Printf("未找到用户名")
		c.JSON(401, gin.H{
			"code":    401,
			"message": "未找到用户信息",
		})
		return
	}

	// 获取查询参数中的开始和结束时间字符串
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	// 定义正则表达式来检测时间字符串是否包含时间部分
	timeRegex := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2})(\.\d+)?Z$`)
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

	// 初始化时间变量
	var startTime, endTime time.Time
	var err error

	// 检查 start_time 是否包含时间部分
	if timeRegex.MatchString(startTimeStr) { // 含时间部分
		startTime, err = time.Parse("2006-01-02T15:04:05Z", startTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "无效的开始时间格式", "error": err.Error()})
			return
		}
	} else if dateRegex.MatchString(startTimeStr) { // 不含时间部分
		startTime, err = time.Parse("2006-01-02", startTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "无效的开始时间格式", "error": err.Error()})
			return
		}
		// 设置开始时间为该日期的零点
		startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, time.UTC)
	} else { // 都不符合
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的开始时间格式"})
		return
	}

	// 检查 end_time 是否包含时间部分
	if timeRegex.MatchString(endTimeStr) { // 含时间部分
		endTime, err = time.Parse("2006-01-02T15:04:05Z", endTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "无效的结束时间格式", "error": err.Error()})
			return
		}
	} else if dateRegex.MatchString(endTimeStr) { // 不含时间部分
		endTime, err = time.Parse("2006-01-02", endTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "无效的结束时间格式", "error": err.Error()})
			return
		}
		// 设置结束时间为该日期的最后一秒
		endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 999999999, time.UTC)
	} else { // 都不符合
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的结束时间格式"})
		return
	}

	// 调用函数获取指定时间范围内的图像信息列表
	images, err := image_r.GetImagesInfoWithinTimeRange(d.DB, username.(string), startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询图像列表失败", "error": err.Error()})
		return
	}

	// 返回图像列表
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询图像列表成功",
		"images":  images,
	})
}

// 获取所有图像信息
// @Summary 获取所有图像信息
// @Description 获取系统中所有图像的信息
// @Tags images
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "获取图像列表成功"
// @Failure 500 {object} map[string]interface{} "获取图像列表失败"
// @Router /image/all [get]
// GetAllImages 获取所有图片信息
// 该函数从数据库中检索所有图片的信息，并以JSON格式返回给客户端
// 参数:
//
//	c *gin.Context: Gin框架的上下文对象，用于处理HTTP请求和响应
func GetAllImages(c *gin.Context) {
	// 调用业务逻辑层函数获取所有图片信息
	images, err := image_r.GetAllImagesInfo(d.DB)
	if err != nil {
		// 如果发生错误，返回500错误，表明服务器内部错误
		c.JSON(http.StatusInternalServerError, gin.H{"message": "获取图像列表失败", "error": err.Error()})
		return
	}

	// 如果成功获取图片信息，返回200状态码和图片列表
	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"message": "获取图像列表成功",
		"images":  images,
	})
}

// GetAllImagesWithLike 获取用户喜欢的图像列表
// 该函数从上下文中获取用户名，然后查询数据库以获取图像信息及用户是否喜欢的状态
// 参数:
//
//	c *gin.Context - Gin框架的上下文对象，用于处理HTTP请求和响应
func GetAllImagesWithLike(c *gin.Context) {
	// 尝试从上下文中获取用户名，如果不存在或为空，则返回401错误
	usernames, exist := c.Get("username")
	if !exist || usernames == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "名字解析出错"})
		return
	}
	// 将获取到的用户名转换为字符串类型
	username, _ := usernames.(string) //当前用户的用户名

	// 调用业务逻辑层，获取所有图像信息及用户喜欢的状态
	images, err := image_r.GetAllImagesInfoWithLikeStatus(d.DB, username)
	if err != nil {
		// 如果发生错误，返回500错误
		c.JSON(http.StatusInternalServerError, gin.H{"message": "获取图像列表失败", "error": err.Error()})
		return
	}
	// 返回成功响应，包含图像列表
	// 改返回 imageResponse 而非 ImageInformation
	c.JSON(http.StatusOK, gin.H{
		"message": "获取图像列表成功",
		"images":  images,
	})
}
