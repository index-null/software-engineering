package user_r

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	d "text-to-picture/models/init"
	u "text-to-picture/models/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetUserById(db *gorm.DB, id int) (*u.UserInformation, error) {
	var user u.UserInformation
	err := db.Table("userinformation").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// 根据用户名查询用户信息
func GetUserByName(db *gorm.DB, username string) (*u.UserInformation, error) {
	var user u.UserInformation
	err := db.Table("userinformation").Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, err
}

// 根据电子邮件查询用户信息
func GetUserByEmail(db *gorm.DB, email string) (*u.UserInformation, error) {
	var user u.UserInformation
	err := db.Table("userinformation").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// 获取所有用户信息并id升序
func GetAllUsers(db *gorm.DB) ([]u.UserInformation, error) {
	var users []u.UserInformation
	result := db.Table("userinformation").Order("id ASC").Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("查询用户列表时发生错误: %v", result.Error)
	}
	return users, nil
}

//-----------------------------------------------------------------------------------------------
func GetUserInfo(c *gin.Context) {

	username := c.Query("username") // 从查询参数中获取用户名
	useremail := c.Query("email")
	userId := c.Query("id")
	userid, err1 := strconv.Atoi(userId)

	var user *u.UserInformation
	var err error
	if  err1 == nil {
		user, err = GetUserById(d.DB, userid)
	}else if username != "" {
		user, err = GetUserByName(d.DB, username)
	} else if useremail != "" {
		user, err = GetUserByEmail(d.DB, useremail)
	}else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
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

func GetAllUsersInfo(c *gin.Context){
	users, err := GetAllUsers(d.DB)
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
