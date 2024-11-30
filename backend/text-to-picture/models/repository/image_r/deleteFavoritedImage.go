package image_r

import (
	"errors"
	"fmt"
	i "text-to-picture/models/image"
	"time"

	"gorm.io/gorm"
)

// 取消图像收藏
func DeleteFavoritedImage(db *gorm.DB, userName string, imageUrl string, create_time time.Time) error {
	var favoritedImage i.FavoritedImages

	// 首先根据用户名和图像URL查找收藏记录
	if err := db.Table("favoritedimage").Where("username = ? AND result = ?", userName, imageUrl).First(&favoritedImage).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("未找到对应的收藏记录")
		}
		return fmt.Errorf("查询收藏记录失败: %v", err)
	}

	// 使用查找到的记录的主键进行删除
	if err := db.Delete(&favoritedImage).Error; err != nil {
		return fmt.Errorf("取消图像收藏失败: %v", err)
	}

	return nil
}