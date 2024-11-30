package like

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"text-to-picture/models/image"
	db "text-to-picture/models/init"
)

type ImageLike struct {
	ID         int
	Picture    string
	Username   string
	Num        int
	CreateTime string
}

func LikeImage(c *gin.Context) {
	// 解析请求中的图片 URL 和 token
	imageURL := c.Query("url")

	if imageURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "Missing image URL"})
		return
	}

	usernames, _ := c.Get("username")
	if usernames == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":  401,
			"error": "名字解析出错"})
		return
	}
	username, _ := usernames.(string)
	tx := db.DB.Begin()
	if tx == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "点赞数据库开始出错"})
		return
	}
	defer func() {
		if tx == nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var imageLike ImageLike

	// 查询用户是否有点赞记录
	if err := tx.Where("username = ? AND picture = ?", username, imageURL).First(&imageLike).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"code":  409,
			"error": "User has already liked this image"})
		return
	}

	// 获取当前点赞数
	var currentLikeCount int
	if err := tx.Model(&image.ImageInformation{}).Where("result = ?", imageURL).Select("likecount").Row().Scan(&currentLikeCount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"error": err.Error()})
		return
	}

	newImageLike := ImageLike{
		Picture:  imageURL,
		Username: username,
		Num:      1,
	}

	if err := tx.Create(&newImageLike).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"error": err.Error()})
		return
	}
	if err := tx.Model(&image.ImageInformation{}).Where("picture = ?", imageURL).Update("likecount", currentLikeCount+1).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{

		"message": "Image liked successfully"})
}
