package user_r

import (
	"fmt"
	u "text-to-picture/models/user"

	"gorm.io/gorm"
)

func DeleteUserByUsername(db *gorm.DB, username string) error {
	if err := db.Table("userinformation").Where("username = ?", username).Delete(&u.UserInformation{}).Error; err != nil {
		return fmt.Errorf("删除用户 %v 失败: %v", username, err)
	}

	return nil
}


