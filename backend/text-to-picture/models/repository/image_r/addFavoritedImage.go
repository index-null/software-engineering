package image_r

import (
	i "text-to-picture/models/image"
	"time"

	"gorm.io/gorm"
)

// 添加收藏的图像
func AddFavoritedImage(db *gorm.DB, userName string, imageUrl string, create_time time.Time) error {
	favoritedImage := i.FavoritedImages{
		UserName:    userName,
		Picture:     imageUrl,
		Create_time: create_time,
	}

	return db.Create(&favoritedImage).Error
}
