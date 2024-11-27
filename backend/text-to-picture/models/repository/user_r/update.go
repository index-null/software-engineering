package user_r

import (
	"errors"
	"fmt"
	u "text-to-picture/models/user"

	"gorm.io/gorm"
)

// 更新用户信息
func UpdateUserInfo(db *gorm.DB, username string, updates map[string]interface{}) error {
	// 查询用户是否存在
	var user u.UserInformation
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return fmt.Errorf("查询用户时发生错误: %v", result.Error)
	}

	// 创建一个临时结构体来存储需要更新的字段
	type UpdateStruct struct {
		Email     *string
		Username  *string
		Password  *string
		AvatarURL *string
		Token     *string
	}

	updateStruct := UpdateStruct{}

	// 反射设置需要更新的字段
	for key, value := range updates {
		switch key {
		case "email":
			if str, ok := value.(string); ok {
				updateStruct.Email = &str
			}
		case "username":
			if str, ok := value.(string); ok {
				updateStruct.Username = &str
			}
		case "password":
			if str, ok := value.(string); ok {
				updateStruct.Password = &str
			}
		case "avatar_url":
			if str, ok := value.(string); ok {
				updateStruct.AvatarURL = &str
			}
		case "token":
			if str, ok := value.(string); ok {
				updateStruct.Token = &str
			}
		}
	}

	// 验证更新的字段
	if updateStruct.Username != nil  {
		return fmt.Errorf("用户名不可修改")
	}
	if updateStruct.Email != nil && *updateStruct.Email == "" {
		return fmt.Errorf("邮箱为空")
	}
	if updateStruct.Password != nil && len(*updateStruct.Password) < 6 {
		return fmt.Errorf("密码少于6位")
	}
	if updateStruct.Email != nil && !isValidEmail(*updateStruct.Email) {
		return fmt.Errorf("邮箱格式不正确")
	}

	// 使用 Updates 方法更新指定字段
	if err := db.Model(&user).Updates(updateStruct).Error; err != nil {
		return fmt.Errorf("更新用户信息失败: %v", err)
	}

	return nil
}
