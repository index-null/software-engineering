package image

import (
	"net/http"
	"strconv"
	services "text-to-picture/services/image_s" // 引入services模块，处理核心业务逻辑

	"github.com/gin-gonic/gin"
)

// 获取图片广场中的图片列表
func GetImages(c *gin.Context) {
	// 从上下文中获取用户名（JWT中提取的用户信息）
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户未认证"})
		return
	}

	// 从请求中获取分页参数
	page := c.DefaultQuery("page", "1")
	size := c.DefaultQuery("size", "10")
	offset, _ := strconv.Atoi(page)
	limit, _ := strconv.Atoi(size)

	// 调用服务层的 GetImages 函数获取图片列表
	images, err := services.GetImagesByUsername(username.(string), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "获取图像列表失败", "error": err})
		return
	}

	// 返回获取到的图片列表
	c.JSON(http.StatusOK, gin.H{"message": "获取图像列表成功", "images": images})
}

func LikeImage(c *gin.Context) {
	// 从上下文中获取用户名（JWT中提取的用户信息）
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户未认证"})
		return
	}

	// 获取请求中的图片ID
	imageID := c.Param("image_id")

	// 调用服务层的 LikeImage 函数处理点赞逻辑
	err := services.LikeImage(imageID, username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "点赞失败", "error": err})
		return
	}

	// 点赞成功，返回成功消息
	c.JSON(http.StatusOK, gin.H{"message": "点赞成功"})
}
