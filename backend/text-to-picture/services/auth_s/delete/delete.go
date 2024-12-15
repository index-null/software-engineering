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

func DeleteUserByName(c *gin.Context) {
	// 从上下文中获取用户名
	userName, exists := c.Get("username")
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

	// 非root用户，不能删除其他某个用户
	if (isOwn != "true" && isOwn != "True") && userName.(string) != "root" {
		c.JSON(http.StatusBadRequest, gin.H{
			//"success": false,
			"message": "非root用户，不可删除其他某个用户",
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
	if isOwn == "true" || isOwn == "True" {
		c.JSON(200,gin.H{"message": username + "的账号注销成功"})
		return 
	}
	c.JSON(200, gin.H{"message": "成功删除用户：" + username})
}
