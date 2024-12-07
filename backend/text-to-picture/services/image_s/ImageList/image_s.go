package ImageList

import (
	"fmt"
	"text-to-picture/models/image"
	db "text-to-picture/models/init"
)

// GetImagesByUsername 获取用户的图片列表
func GetImagesByUsername(username string, offset, limit int) ([]image.ImageInformation, error) {
	var images []image.ImageInformation
	// 查询数据库，根据用户名获取图片列表
	err := db.DB.Preload("User").Where("username = ?", username).Limit(limit).Offset(offset * limit).Find(&images).Error
	if err != nil {
		return nil, fmt.Errorf("查询图片失败: %v", err)
	}
	return images, nil
}
