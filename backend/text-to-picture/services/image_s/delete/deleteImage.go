package delete

import (
	//"fmt"
	"fmt"
	"log"
	"net/http"

	//"strconv"
	//"strings"

	d "text-to-picture/models/init"
	"text-to-picture/models/repository/image_r"
	"text-to-picture/models/repository/user_r"

	u "text-to-picture/models/user"

	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	ImageUrl string `json:"url"`
	Id       int    `json:"id" `
}

type BatchDeleteRequestBody struct {
	Urls []string `json:"urls"`
	Ids	 []int	  `json:"ids"`
}

// 删除用户的单个图像
func DeleteUserOneImage(c *gin.Context) {
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

	var err error
	var requestBody RequestBody
	// 解析请求体
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的请求格式", "error": err.Error()})
		return
	}
	
	// 开始事务
	tx := d.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "开始事务失败", "error": tx.Error.Error()})
		return
	}

	err = image_r.DeleteUserOneImage(tx, requestBody.ImageUrl,username,requestBody.Id)
	if err != nil {
		tx.Rollback() // 回滚事务
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除用户的一张图像失败", "error": err.Error()})
		return
	}

	if err = tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "提交事务失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "成功删除用户的一张图像"})
}

// 删除用户多张图像
func DeleteUserImagesBatch(c *gin.Context) {
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
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的请求格式", "error": err.Error()})
		return
	}

	// 开始事务
	tx := d.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "开始事务失败", "error": tx.Error.Error()})
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器内部错误"})
			return
		}
	}()

	var errors []string
	
	if len(requestBody.Urls) >0 && len(requestBody.Ids) == 0{		
		for _, url := range requestBody.Urls {
			err := image_r.DeleteUserOneImage(tx, url, username, 0)
			if err != nil {
				errors = append(errors, fmt.Sprintf("删除URL %s 失败：%v", url,err))
			}
		}
	}else if len(requestBody.Ids)>0 && len(requestBody.Urls)==0{
		for _, id := range requestBody.Ids{
			err :=image_r.DeleteUserOneImage(tx, "", username, id)
			if err != nil {
				errors = append(errors, fmt.Sprintf("删除ID %d 失败：%v",id,err))
			}
		}
	}else{
		tx.Rollback()
		c.JSON(http.StatusBadRequest,gin.H{"message":"请提供有效的urls或ids列表，但不要同时提供这两个列表"})
		return
	}

	if len(errors)>0{
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,gin.H{
			"message":"部分或全部图像删除失败，撤销删除",
			"errors":errors,
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "提交事务失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "成功删除用户指定的图像"})
}


// 删除单个图像
func DeleteOneImage(c *gin.Context) {
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

	// 不是root用户，不能删除图像
	if username != "root" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "非root用户，不可删除某张图像",
		})
		return
	}

	url := c.Query("url")

	// 开始事务
	tx := d.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "开始事务失败", "error": tx.Error.Error()})
		return
	}

	err := image_r.DeleteOneImage(tx, url)
	if err != nil {
		tx.Rollback() // 回滚事务
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除图像失败", "error": err.Error()})
		return
	}

	if err = tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "提交事务失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "成功删除图像"})
}

// 删除指定用户的全部图像
func DeleteUserImages(c *gin.Context) {
	userName, exists := c.Get("username")
	if !exists {
		log.Printf("未找到用户名")
		c.JSON(401, gin.H{
			"success": false,
			"message": "未找到用户信息",
		})
		return
	}

	// 不是root用户，不能删除图像
	if userName.(string) != "root" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "非root用户，不可删除某个用户的所有图像",
		})
		return
	}

	username := c.Query("username")

	// 开始事务
	tx := d.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "开始事务失败", "error": tx.Error.Error()})
		return
	}

	var user *u.UserInformation
	user, err := user_r.GetUserByName(tx, username)
	if err != nil {
		tx.Rollback() // 回滚事务
		c.JSON(http.StatusNotFound, gin.H{"message": "用户" + username + "不存在", "error": err.Error()})
		return
	}

	err = image_r.DeleteUserAllImages(tx, user.UserName)
	if err != nil {
		tx.Rollback() // 回滚事务
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除用户" + username + "的所有图像失败", "error": err.Error()})
		return
	}

	if err = tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "提交事务失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "成功删除用户" + username + "的所有图像"})
}

// 删除全部图像
func DeleteAllImages(c *gin.Context) {
	userName, exists := c.Get("username")
	if !exists {
		log.Printf("未找到用户名")
		c.JSON(401, gin.H{
			"success": false,
			"message": "未找到用户信息",
		})
		return
	}

	// 不是root用户，不能删除图像
	if userName.(string) != "root" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "非root用户，不可删除全部图像",
		})
		return
	}

	// 开始事务
	tx := d.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "开始事务失败", "error": tx.Error.Error()})
		return
	}

	err := image_r.DeleteAllImages(tx)
	if err != nil {
		tx.Rollback() // 回滚事务
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除全部图像失败", "error": err.Error()})
		return
	}

	if err = tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "提交事务失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "成功删除全部图像"})
}
