package delete

import (
	"fmt"
	"log"
	"net/http"

	d "text-to-picture/models/init"
	u "text-to-picture/models/repository/user_r"

	"github.com/gin-gonic/gin"
	//"gorm.io/gorm"
)

// @Summary 删除用户
// @Description 根据用户名删除用户，只有root用户才能删除其他用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param username query string true "用户名"
// @Success 200 {object} map[string]interface{} "成功删除用户"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 404 {object} map[string]interface{} "用户不存在"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /auth/root/deleteOneUser [delete]
func DeleteUserByName(c *gin.Context) {
	// 从上下文中获取用户名
	userName, exists := c.Get("username") // 从上下文中获取当前登录用户的用户名
	if !exists {
		log.Printf("未找到用户名")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未找到用户信息",
		})
		return
	}

	fmt.Println("当前的登录用户为：" + userName.(string))

	isOwn := c.Query("isOwn")

	// 判断是不是既非root用户也非账号注销这两种情况
	if (isOwn != "true" && isOwn != "True") && userName.(string) != "root" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "非root用户，不可删除其他某个用户",
		})
		return
	}

	// isOwn为true则表示用户的账号注销操作，否则为root用户的删除用户操作
	var username string
	if isOwn != "true" {
		username = c.Query("username")
	} else if isOwn == "true" {
		username = userName.(string)
	}

	// 开始事务
	tx := d.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "开始事务失败", "error": tx.Error.Error()})
		return
	}

	// 检查用户是否存在
	exist, err := u.IsExist(tx, username)
	if err != nil {
		tx.Rollback() // 回滚事务
		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询用户失败", "error": err.Error()})
		return
	}

	if !exist { // 不存在
		tx.Rollback() // 回滚事务
		c.JSON(http.StatusNotFound, gin.H{"message": "用户不存在"})
		return
	}

	// 删除用户
	err = u.DeleteUserByUsername(tx, username)
	if err != nil {
		tx.Rollback() // 回滚事务
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除用户失败", "error": err.Error()})
		return
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "提交事务失败", "error": err.Error()})
		return
	}

	// 返回结果
	if isOwn == "true" || isOwn == "True" {
		c.JSON(200, gin.H{"message": username + "的账号注销成功"})
		return
	}
	c.JSON(200, gin.H{"message": "成功删除用户：" + username})
}
