package user_r

import (
	"fmt"
	u "text-to-picture/models/user"
	"gorm.io/gorm"
)

func IsExist(db *gorm.DB, username string) (bool, error) {
	var count int64
	err := db.Table("userinformation").Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count <= 0 {
		return false, nil
	} else {
		return true, nil
	}
}

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
