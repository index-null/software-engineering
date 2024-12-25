package user_r

import (
	"fmt"
	u "text-to-picture/models/user"

	"gorm.io/gorm"
)

func DeleteUserByUsername(db *gorm.DB, username string) error {
	// 在 userinformation 表中删除指定用户名的用户
	if err := db.Table("userinformation").Where("username = ?", username).Delete(&u.UserInformation{}).Error; err != nil {
		return fmt.Errorf("删除用户 %v 失败: %v", username, err) // 返回删除失败的错误信息
	}

	return nil // 删除成功
}


