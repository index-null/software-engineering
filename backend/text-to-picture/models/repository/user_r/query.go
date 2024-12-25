package user_r

import (
	"fmt"
	u "text-to-picture/models/user"
	"gorm.io/gorm"
)

// IsExist 检查用户是否存在
func IsExist(db *gorm.DB, username string) (bool, error) {
	var count int64
	// 查询指定用户名的记录数
	err := db.Table("userinformation").Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, err // 返回查询错误
	}
	if count <= 0 {
		return false, nil // 用户不存在
	} else {
		return true, nil // 用户存在
	}
}

func GetUserById(db *gorm.DB, id int) (*u.UserInformation, error) {
	var user u.UserInformation
	err := db.Table("userinformation").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err // 返回查询错误
	}
	return &user, nil // 返回用户信息
}

// 根据用户名查询用户信息
func GetUserByName(db *gorm.DB, username string) (*u.UserInformation, error) {
	var user u.UserInformation
	err := db.Table("userinformation").Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err // 返回查询错误
	}
	return &user, err // 返回用户信息
}

// 根据电子邮件查询用户信息
func GetUserByEmail(db *gorm.DB, email string) (*u.UserInformation, error) {
	var user u.UserInformation
	err := db.Table("userinformation").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err // 返回查询错误
	}
	return &user, nil // 返回用户信息
}

// 获取所有用户信息并id升序
func GetAllUsers(db *gorm.DB) ([]u.UserInformation, error) {
	var users []u.UserInformation
	result := db.Table("userinformation").Order("id ASC").Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("查询用户列表时发生错误: %v", result.Error) // 返回查询错误
	}
	return users, nil // 返回用户信息
}
