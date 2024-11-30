package update

import (
	"log"
	"net/http"
	d "text-to-picture/models/init"
	"text-to-picture/models/repository/user_r"

	"github.com/gin-gonic/gin"
)

// 更新用户信息
func UpdateUser(c *gin.Context) {//不能更新用户名
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

	// 获取用户名
	//username := c.Param("username")

	// 定义用于接收 JSON 数据的结构体
	var input map[string]interface{}

	// 解析 JSON 数据
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求数据格式错误","error":err})
		return
	}

	// 更新用户信息
	if err := user_r.UpdateUserInfo(d.DB, username.(string), input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "更新用户信息失败", "error": err.Error()})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"message": "用户信息更新成功",
	})
}