package user_r

import (
	"errors"
	"fmt"
	"log"
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
	if updateStruct.Username != nil {
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

	fmt.Printf("准备更新的用户信息: %+v\n", updateStruct)

	tx := db.Begin()
    if err := tx.Error; err != nil {
        return fmt.Errorf("开始事务失败: %v", err)
    }
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            fmt.Printf("事务处理过程中出现panic: %v\n", r)
        }
    }()

	log.Printf("正在尝试更新用户 %s 的信息...", username)
	
	
    // 执行数据库操作
     // 执行数据库操作
	result = tx.Model(&user).Updates(updateStruct)
	if result.Error != nil {
		tx.Rollback()
		log.Printf("更新用户信息失败: %v", result.Error)
		return fmt.Errorf("更新用户信息失败: %v", result.Error)
	}
	log.Printf("成功更新了用户 %s 的信息，受影响行数: %d", username, result.RowsAffected)

	// 提交事务前打印日志
	log.Printf("正在提交事务...")

	if err := tx.Commit().Error; err != nil {
		log.Printf("提交事务失败: %v", err)
		return fmt.Errorf("提交事务失败: %v", err)
	}

	log.Printf("事务提交成功")
	return nil
}
