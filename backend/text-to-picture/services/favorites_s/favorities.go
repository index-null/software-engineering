package favorites_s

import (
	"log"
	"net/http"
	"strconv"
	i "text-to-picture/models/image"

	d "text-to-picture/models/init"
	"text-to-picture/models/repository/image_r"

	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	ImageUrl string `json:"url"`
	Id       int    `json:"id" `
}

// 收藏图像
// @Summary 收藏图像
// @Description 根据图像URL或ID收藏图像
// @Tags favorites
// @Accept json
// @Produce json
// @Param requestBody body favorites_s.RequestBody true "请求体"
// @Success 200  {object} map[string]interface{} "图像收藏成功"
// @Failure 400  {object} map[string]interface{} "无效的请求格式"
// @Failure 401  {object} map[string]interface{} "未找到用户信息"
// @Failure 404  {object} map[string]interface{} "未找到对应的图像"
// @Failure 409  {object} map[string]interface{} "该图像已经被收藏过"
// @Failure 500  {object} map[string]interface{} "检查收藏状态失败"
// @Router /auth/addFavoritedImage [post]
func AddFavoritedImage(c *gin.Context) {
	var requestBody RequestBody
	var imageInfo *i.ImageInformation
	var err error
	// 解析请求体
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的请求格式", "error": err.Error()})
		return
	}

	// 查询图像信息
	if requestBody.ImageUrl != "" {
		imageInfo, err = image_r.GetImageByUrl(d.DB, requestBody.ImageUrl)
	} else if requestBody.Id > 0 {
		imageInfo, err = image_r.GetImageById(d.DB, requestBody.Id)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无有效的图像id或url", "error": "id 必须大于 0 或者 url 不得为空"})
		return
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "未找到对应的图像", "error": err.Error()})
		return
	}

	// 从上下文中获取用户名
	username, exists := c.Get("username")
	if !exists {
		log.Printf("未找到用户名")
		c.JSON(401, gin.H{
			"success": false,
			"message": "未找到用户信息",
		})
		return
	}

	// 在添加收藏之前检查用户是否已收藏该图像
	isFavorited, err := image_r.IsImageFavoritedByUser(d.DB, username.(string), imageInfo.Picture)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "检查收藏状态失败", "error": err.Error()})
		return
	}

	if isFavorited {
		c.JSON(http.StatusConflict, gin.H{"message": "该图像已经被收藏过"})
		return
	}

	// 添加收藏图像
	err = image_r.AddFavoritedImage(d.DB, username.(string), imageInfo.Picture, imageInfo.Create_time)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "收藏图像失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "图像收藏成功"})

}

// @Summary 取消图像收藏
// @Description 根据图像URL或ID取消收藏图像
// @Tags favorites
// @Accept json
// @Produce json
// @Param url query string false "图像URL"
// @Param id query int false "图像ID"
// @Success 200  {object} map[string]interface{} "取消图像收藏成功"
// @Failure 400  {object} map[string]interface{} "无效的请求格式"
// @Failure 401  {object} map[string]interface{} "未找到用户信息"
// @Failure 404  {object} map[string]interface{} "未找到对应的图像"
// @Failure 409  {object} map[string]interface{} "该图像未被收藏过，不可取消收藏"
// @Failure 500  {object} map[string]interface{} "检查收藏状态失败"
// @Router /auth/deleteFavoritedImage [delete]
// 取消图像收藏 DELETE方法
func DeleteFavoritedImage(c *gin.Context) {
	var imageInfo *i.ImageInformation
	var err error
	url := c.Query("url")
	idStr := c.Query("id")

	// 查询图像信息
	if url != "" {
		imageInfo, err = image_r.GetImageByUrl(d.DB, url)

	} else if idStr != "" {
			id, err1 := strconv.Atoi(idStr)
			if err1 != nil || id <= 0 {
				c.JSON(http.StatusBadRequest, gin.H{"message": "无有效的图像id或url", "error": "id 必须大于 0 或者 url 不得为空"})
				return
			}
			var imageInfo1 *i.FavoritedImages
			imageInfo1, err = image_r.GetFavoritedImageById(d.DB, imageInfo.ID)
			imageInfo.Picture = imageInfo1.Picture

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无有效的图像id或url", "error": "id 必须大于 0 或者 url 不得为空"})
		return
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "未找到对应的图像", "error": err.Error()})
		return
	}

	// 从上下文中获取用户名
	username, exists := c.Get("username")
	if !exists {
		log.Printf("未找到用户名")
		c.JSON(401, gin.H{
			"success": false,
			"message": "未找到用户信息",
		})
		return
	}

	// 在取消收藏之前检查用户是否已收藏该图像
	var isFavorited bool
	isFavorited, err = image_r.IsImageFavoritedByUser(d.DB, username.(string), imageInfo.Picture)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "检查收藏状态失败", "error": err.Error()})
		return
	}

	if !isFavorited {
		c.JSON(http.StatusConflict, gin.H{"message": "该图像未被收藏过，不可取消收藏"})
		return
	}

	// 取消收藏
	err = image_r.DeleteFavoritedImage(d.DB, username.(string), imageInfo.Picture)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "取消图像收藏失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "取消图像收藏成功"})

}
