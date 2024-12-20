package image_r

import (
	"errors"
	"fmt"
	i "text-to-picture/models/image"

	"gorm.io/gorm"
)

// 删除指定一张图像
func DeleteUserOneImage(db *gorm.DB, imageUrl string, username string, id int) error {
	var image i.ImageInformation

	if imageUrl != "" {
		// 首先根据用户名和图像URL查找图像记录
		if err := db.Table("imageinformation").Where("picture = ? AND username = ?", imageUrl, username).First(&image).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("数据库中未找到对应的图像")
			}
			return fmt.Errorf("查询图像失败: %v", err)
		}
	}else if id > 0{
		// 首先根据用户名和图像URL查找图像记录
		if err := db.Table("imageinformation").Where("id = ? AND username = ?", id, username).First(&image).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("数据库中未找到对应的图像")
			}
			return fmt.Errorf("查询图像失败: %v", err)
		}
	}else{
		return fmt.Errorf("没有有效的url或id")
	}

	// 使用查找到的记录的主键进行删除
	if err := db.Delete(&image).Error; err != nil {
		return fmt.Errorf("删除用户的一张图像失败: %v", err)
	}

	return nil
}

// 删除指定一张图像
func DeleteOneImage(db *gorm.DB, imageUrl string) error {
	var image i.ImageInformation

	// 首先根据用户名和图像URL查找图像记录
	if err := db.Table("imageinformation").Where("picture = ?", imageUrl).First(&image).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("数据库中未找到对应的图像")
		}
		return fmt.Errorf("查询图像失败: %v", err)
	}

	// 使用查找到的记录的主键进行删除
	if err := db.Delete(&image).Error; err != nil {
		return fmt.Errorf("删除图像失败: %v", err)
	}

	return nil
}

// 批量删除指定用户的多个图像
func DeleteUserSomeImages(db *gorm.DB, userName string, imageUrls []string) error {
	if len(imageUrls) == 0 {
		return fmt.Errorf("没有提供要删除的图像URL")
	}

	err := db.Table("imageinformation").Where("username = ? AND picture IN (?)", userName, imageUrls).Delete(&i.ImageInformation{}).Error

	if err != nil {
		return fmt.Errorf("删除用户 %s 的多个图像失败：%v", userName, err)
	}

	return nil
}

// 删除用户的所有图像
func DeleteUserAllImages(db *gorm.DB, userName string) error {

	if err := db.Where("username = ?", userName).Delete(&i.ImageInformation{}).Error; err != nil {
		return fmt.Errorf("删除用户 %s 的所有图像失败: %v", userName, err)
	}

	return nil
}

// 删除所有图像
func DeleteAllImages(db *gorm.DB) error {
	if err := db.Table("imageinformation").Delete(&i.ImageInformation{}).Error; err != nil {
		return fmt.Errorf("删除 imageinformation 表中的所有记录失败: %v", err)
	}

	return nil
}
