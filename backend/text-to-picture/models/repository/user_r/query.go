package repository

import (
	d "gocode/backend/backend/text-to-picture/models/init"
	u "gocode/backend/backend/text-to-picture/models/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 根据用户名查询用户信息
func GetUserByUserName(db *gorm.DB, username string) (*u.UserLogin, error) {
	var user u.UserLogin
	err := db.Where("user_name = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// 根据电子邮件查询用户信息
func GetUserByEmail(db *gorm.DB, email string) (*u.UserLogin, error) {
	var user u.UserLogin
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserInfo(c *gin.Context) {

	username := c.Query("user_name") // 从查询参数中获取用户名
	useremail := c.Query("user_email")

	var user *u.UserLogin
	var err error
	if username != "" {
		user, err = GetUserByUserName(d.DB, username)
	} else if useremail != "" {
		user, err = GetUserByEmail(d.DB, username)
	}
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "用户未找到"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, user)
}
