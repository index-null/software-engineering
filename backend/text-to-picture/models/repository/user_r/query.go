package user_r

import (
	"errors"
	d "text-to-picture/models/init"
	u "text-to-picture/models/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 根据用户名查询用户信息
func GetUserByName(db *gorm.DB, username string) (*u.Login, error) {
	var user u.Login
	err := db.Where("user_name = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// 根据电子邮件查询用户信息
func GetUserByEmail(db *gorm.DB, email string) (*u.Login, error) {
	var user u.Login
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserInfo(c *gin.Context) {

	username := c.Query("user_name") // 从查询参数中获取用户名
	useremail := c.Query("user_email")

	var user *u.Login
	var err error
	if username != "" {
		user, err = GetUserByName(d.DB, username)
	} else if useremail != "" {
		user, err = GetUserByEmail(d.DB, username)
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "用户未找到"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, user)
}
