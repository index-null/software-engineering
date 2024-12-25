package like

import (
	"errors"
	"net/http"
	"text-to-picture/models/image"
	db "text-to-picture/models/init"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReqBody struct {
	URL string `json:"url"` // 图片的 URL
}

// @Summary 点赞图片
// @Description 点赞图片接口
// @Tags image
// @Accept json
// @Produce json
// @Param requestBody body like.ReqBody true "图片 URL"
// @Success 200 {object} map[string]interface{} "点赞成功"
// @Failure 400 {object} map[string]interface{} "缺少图片 URL"
// @Failure 401 {object} map[string]interface{} "名字解析出错"
// @Failure 409 {object} map[string]interface{} "用户已点赞该图片"
// @Failure 500 {object} map[string]interface{} "数据库操作错误"
// @Router /auth/like [put]
func LikeImage(c *gin.Context) {
	// 解析请求中的图片 URL 和 token
	var reqBody ReqBody
	err := c.BindJSON(&reqBody)
	if err != nil {
		// 请求体解析失败
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "请求体解析失败",
		})
		return
	}

	imageURL := reqBody.URL
	if imageURL == "" {
		// 如果未提供图片 URL
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "Missing image URL",
		})
		return
	}

	// 从上下文中获取用户名
	usernames, exist := c.Get("username")
	if !exist || usernames == "" {
		// 如果用户名不存在或为空
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":  401,
			"error": "名字解析出错",
		})
		return
	}
	username, _ := usernames.(string)

	// 开始事务
	tx := db.DB.Begin()
	if tx.Error != nil {
		// 事务启动失败
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"error": "点赞数据库开始出错",
		})
		return
	}

	defer func() {
		// 确保事务在发生 panic 或其他错误时回滚
		if r := recover(); r != nil {
			tx.Rollback()
		} else if tx.Error != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var imageLike image.ImageLike

	// 检查用户是否已对图片点赞
	if err := tx.Where("username = ? AND picture = ?", username, imageURL).First(&imageLike).Error; err == nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果用户已点赞，返回冲突状态码
		c.JSON(http.StatusConflict, gin.H{
			"code":  409,
			"error": "用户已经点赞过该图片",
		})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 发生其他非“记录未找到”的错误
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"error": "Database query error",
		})
		return
	}

	// 获取当前点赞数
	var currentLikeCount int
	if err := tx.Model(&image.ImageInformation{}).Where("picture = ?", imageURL).Select("likecount").Row().Scan(&currentLikeCount); err != nil {
		// 查询点赞数失败
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"error": err.Error(),
		})
		return
	}

	// 创建新的点赞记录
	newImageLike := image.ImageLike{
		Picture:  imageURL,
		UserName: username,
		Num:      1,
	}

	if err := tx.Create(&newImageLike).Error; err != nil {
		// 创建点赞记录失败
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"error": err.Error(),
		})
		return
	}

	// 更新图片的点赞数
	if err := tx.Model(&image.ImageInformation{}).Where("picture = ?", imageURL).Update("likecount", currentLikeCount+1).Error; err != nil {
		// 更新点赞数失败
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"error": err.Error(),
		})
		return
	}

	// 返回点赞成功的信息
	c.JSON(http.StatusOK, gin.H{
		"code":          200,
		"current_likes": currentLikeCount + 1,
		"message":       "Image liked successfully",
	})
}
