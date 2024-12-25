package update

import (
	"log"
	"net/http"
	d "text-to-picture/models/init"
	"text-to-picture/models/repository/user_r"

	"github.com/gin-gonic/gin"
)

// 更新用户信息
// @Summary 更新用户信息
// @Description 更新用户的详细信息（不能更新用户名）
// @Tags users
// @Accept json
// @Produce json
// @Param requestBody body map[string]interface{} true "用户信息更新数据"
// @Success 200 {object} map[string]interface{} "用户信息更新成功"
// @Failure 400 {object} map[string]interface{} "请求数据格式错误"
// @Failure 401 {object} map[string]interface{} "未找到用户信息"
// @Failure 500 {object} map[string]interface{} "更新用户信息失败"
// @Router /auth/user/update [put]
func UpdateUser(c *gin.Context) { //不能更新用户名
	// 从上下文中获取用户名
	username, exists := c.Get("username")
	if !exists { // 未找到用户名，说明用户未登录或token有问题
		log.Printf("未找到用户名")
		c.JSON(401, gin.H{
			"success": false,
			"message": "未找到用户信息",
		})
		return
	}

	// 定义用于接收 JSON 数据的结构体
	var input map[string]interface{}

	// 解析 JSON 数据
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求数据格式错误", "error": err})
		return
	}

	// 更新用户信息
	if err := user_r.UpdateUserInfo(d.DB, username.(string), input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新用户信息失败", "error": err.Error()})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "用户信息更新成功",
	})
}
