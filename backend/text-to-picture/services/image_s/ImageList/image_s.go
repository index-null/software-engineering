package ImageList

import (
	"fmt"
	"text-to-picture/models/image"
	db "text-to-picture/models/init"
)

// GetImagesByUsername 获取用户的图片列表
// username: 用户名，用于查询特定用户的图片
// offset: 查询结果的偏移量，用于分页
// limit: 每页返回的图片数量
// 返回值: 图片信息列表和可能的错误信息
func GetImagesByUsername(username string, offset, limit int) ([]image.ImageInformation, error) {
	var images []image.ImageInformation
	// 查询数据库，根据用户名获取图片列表
	// 使用 Preload("User") 预加载关联的用户信息
	// Limit(limit) 和 Offset(offset * limit) 用于分页查询
	err := db.DB.Preload("User").Where("username = ?", username).Limit(limit).Offset(offset * limit).Find(&images).Error
	if err != nil {
		// 如果查询失败，返回错误信息
		return nil, fmt.Errorf("查询图片失败: %v", err)
	}
	// 返回查询到的图片列表和 nil 表示没有错误
	return images, nil
}
