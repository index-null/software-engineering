package delete

import (
	"fmt"
	"log"
	"net/http"

	d "text-to-picture/models/init"
	u "text-to-picture/models/repository/user_r"

	"github.com/gin-gonic/gin"
)

func DeleteUserByName(c *gin.Context) {
	// 从上下文中获取用户名
	userName, exists := c.Get("username")
	fmt.Println("当前的登录用户为：" + userName.(string))
	if !exists {
		log.Printf("未找到用户名")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未找到用户信息",
		})
		return
	}

	// 非root用户，不能删除其他某个用户
	if userName.(string) != "root" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "非root用户，不可删除其他某个用户",
		})
		return
	}

	username := c.Query("username")

	exist, err := u.IsExist(d.DB, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询用户失败", "error": err.Error()})
		return
	}

	if !exist { // 不存在
		c.JSON(http.StatusNotFound, gin.H{"message": "用户不存在"})
		return
	}

	err = u.DeleteUserByUsername(d.DB, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除用户失败", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "成功删除用户" + username})
}
