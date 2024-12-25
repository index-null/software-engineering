package delete

import (
	"fmt"
	"log"
	"net/http"
	d "text-to-picture/models/init"
	"text-to-picture/models/repository/image_r"
	"text-to-picture/models/repository/user_r"
	u "text-to-picture/models/user"

	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	ImageUrl string `json:"url"` // 单张图像的URL
	Id       int    `json:"id"`  // 图像的ID
}

type BatchDeleteRequestBody struct {
	Urls []string `json:"urls"` // 批量删除的图像URL列表
	Ids  []int    `json:"ids"`  // 批量删除的图像ID列表
}

// DeleteUserOneImage 删除用户的单个图像
func DeleteUserOneImage(c *gin.Context) {
	// 从上下文中获取用户名
	userName, exists := c.Get("username")
	if !exists {
		log.Printf("未找到用户名")
		c.JSON(401, gin.H{
			"success": false,
			"message": "未找到用户信息",
		})
		return
	}

	username := userName.(string)

	var requestBody RequestBody
	// 解析请求体中的 JSON 数据
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的请求格式", "error": err.Error()})
		return
	}

	// 开始数据库事务
	tx := d.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "开始事务失败", "error": tx.Error.Error()})
		return
	}

	// 删除用户的一张图像
	err := image_r.DeleteUserOneImage(tx, requestBody.ImageUrl, username, requestBody.Id)
	if err != nil {
		tx.Rollback() // 回滚事务
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除用户的一张图像失败", "error": err.Error()})
		return
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "提交事务失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "成功删除用户的一张图像"})
}

// DeleteUserImagesBatch 删除用户多张图像
func DeleteUserImagesBatch(c *gin.Context) {
	// 从上下文中获取用户名
	userName, exists := c.Get("username")
	if !exists {
		log.Printf("未找到用户名")
		c.JSON(401, gin.H{
			"success": false,
			"message": "未找到用户信息",
		})
		return
	}

	username := userName.(string)

	var requestBody BatchDeleteRequestBody
	// 解析请求体中的 JSON 数据
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的请求格式", "error": err.Error()})
		return
	}

	// 开始数据库事务
	tx := d.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "开始事务失败", "error": tx.Error.Error()})
		return
	}

	// 确保事务在异常时回滚
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器内部错误"})
		}
	}()

	var errors []string

	// 根据请求体中的 Urls 或 Ids 删除图像
	if len(requestBody.Urls) > 0 && len(requestBody.Ids) == 0 {
		for _, url := range requestBody.Urls {
			err := image_r.DeleteUserOneImage(tx, url, username, 0)
			if err != nil {
				errors = append(errors, fmt.Sprintf("删除URL %s 失败：%v", url, err))
			}
		}
	} else if len(requestBody.Ids) > 0 && len(requestBody.Urls) == 0 {
		for _, id := range requestBody.Ids {
			err := image_r.DeleteUserOneImage(tx, "", username, id)
			if err != nil {
				errors = append(errors, fmt.Sprintf("删除ID %d 失败：%v", id, err))
			}
		}
	} else {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": "请提供有效的urls或ids列表，并且不要同时提供这两个列表"})
		return
	}

	// 如果有错误，回滚事务并返回错误信息
	if len(errors) > 0 {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "部分或全部图像删除失败，撤销删除",
			"errors":  errors,
		})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "提交事务失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "成功删除用户指定的图像"})
}

// DeleteOneImage 删除单个图像
func DeleteOneImage(c *gin.Context) {
	// 从上下文中获取用户名
	userName, exists := c.Get("username")
	if !exists {
		log.Printf("未找到用户名")
		c.JSON(401, gin.H{
			"success": false,
			"message": "未找到用户信息",
		})
		return
	}

	username := userName.(string)

	// 检查用户是否为 root
	if username != "root" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "非root用户，不可删除某张图像",
		})
		return
	}

	// 获取 URL 参数中的图像 URL
	url := c.Query("url")

	// 开始数据库事务
	tx := d.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "开始事务失败", "error": tx.Error.Error()})
		return
	}

	// 删除单个图像
	err := image_r.DeleteOneImage(tx, url)
	if err != nil {
		tx.Rollback() // 回滚事务
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除图像失败", "error": err.Error()})
		return
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "提交事务失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "成功删除图像"})
}

// DeleteUserImages 删除指定用户的全部图像
func DeleteUserImages(c *gin.Context) {
	// 从上下文中获取用户名
	userName, exists := c.Get("username")
	if !exists {
		log.Printf("未找到用户名")
		c.JSON(401, gin.H{
			"success": false,
			"message": "未找到用户信息",
		})
		return
	}

	// 检查用户是否为 root
	if userName.(string) != "root" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "非root用户，不可删除某个用户的所有图像",
		})
		return
	}

	// 获取 URL 参数中的目标用户名
	username := c.Query("username")

	// 开始数据库事务
	tx := d.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "开始事务失败", "error": tx.Error.Error()})
		return
	}

	// 检查目标用户是否存在
	var user *u.UserInformation
	user, err := user_r.GetUserByName(tx, username)
	if err != nil {
		tx.Rollback() // 回滚事务
		c.JSON(http.StatusNotFound, gin.H{"message": "用户" + username + "不存在", "error": err.Error()})
		return
	}

	// 删除用户的所有图像
	err = image_r.DeleteUserAllImages(tx, user.UserName)
	if err != nil {
		tx.Rollback() // 回滚事务
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除用户" + username + "的所有图像失败", "error": err.Error()})
		return
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "提交事务失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "成功删除用户" + username + "的所有图像"})
}

// DeleteAllImages 删除全部图像
func DeleteAllImages(c *gin.Context) {
	// 从上下文中获取用户名
	userName, exists := c.Get("username")
	if !exists {
		log.Printf("未找到用户名")
		c.JSON(401, gin.H{
			"success": false,
			"message": "未找到用户信息",
		})
		return
	}

	// 检查用户是否为 root
	if userName.(string) != "root" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "非root用户，不可删除全部图像",
		})
		return
	}

	// 开始数据库事务
	tx := d.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "开始事务失败", "error": tx.Error.Error()})
		return
	}

	// 删除所有图像
	err := image_r.DeleteAllImages(tx)
	if err != nil {
		tx.Rollback() // 回滚事务
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除全部图像失败", "error": err.Error()})
		return
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "提交事务失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "成功删除全部图像"})
}
