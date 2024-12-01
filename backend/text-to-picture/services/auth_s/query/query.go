package query

import (
	"errors"
	"log"
	"net/http"
	//"strconv"
	d "text-to-picture/models/init"
	"text-to-picture/models/repository/user_r"
	u "text-to-picture/models/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary 获取用户信息
// @Description 根据用户名、邮箱或用户ID获取用户信息
// @Tags user
// @Accept json
// @Produce json
// @Param username query string false "用户名"
// @Param email query string false "邮箱"
// @Param id query int false "用户ID"
// @Success 200 {object} map[string]interface{} "获取用户信息成功"
// @Failure 400 {object} map[string]interface{} "无效的请求数据"
// @Failure 404 {object} map[string]interface{} "用户未找到"
// @Failure 500 {object} map[string]interface{} "查询失败"
// @Router /getuserinfo [get]
func GetUserInfo(c *gin.Context) {

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

	// username := c.Query("username") // 从查询参数中获取用户名
	// useremail := c.Query("email")
	// userId := c.Query("id")
	// userid, err1 := strconv.Atoi(userId)

	var user *u.UserInformation
	var err error
	// if err1 == nil {
	// 	user, err = user_r.GetUserById(d.DB, userid)
	// } else if username != "" {
	// 	user, err = user_r.GetUserByName(d.DB, username)
	// } else if useremail != "" {
	// 	user, err = user_r.GetUserByEmail(d.DB, useremail)
	// } else {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
	// 	return
	// }
	user, err = user_r.GetUserByName(d.DB, username.(string))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "用户未找到"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询失败", "error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func GetAllUsersInfo(c *gin.Context) {
	users, err := user_r.GetAllUsers(d.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "获取用户列表失败", "error": err.Error()})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"message": "获取用户列表成功",
		"users":   users,
	})
}
