package image_r

import (
	"fmt"
	"gorm.io/gorm"
	i "text-to-picture/models/image"
	"time"
)

// // 根据用户ID查询相关图片
// func GetUserImagesByUserId(db *gorm.DB, userId int) ([]i.ImageInformation, error) {
// 	var images []i.ImageInformation
// 	err := db.Table("imageinformation").Where("userid = ?", userId).Find(&images).Error // 使用 Find 而不是 First
// 	if err != nil {
// 		return nil, err
// 	}

// 	return images, nil
// }

// // 根据用户ID查询收藏的图片
// func GetUserFavoritedImagesByUserId(db *gorm.DB, userId int) ([]i.ImageInformation, error) {
// 	var images []i.ImageInformation
// 	err := db.Table("favoritedimage").Where("userid = ?", userId).Find(&images).Error // 使用 Find 而不是 First
// 	if err != nil {
// 		return nil, err
// 	}

// 	return images, nil
// }

//-----------------------------------------------获取指定用户的图像

// 根据用户名查询相关图片
func GetUserImagesByUsername(db *gorm.DB, username string) ([]i.ImageInformation, error) {
	var images []i.ImageInformation
	err := db.Table("imageinformation").Where("username = ?", username).Find(&images).Error // 使用 Find 而不是 First
	if err != nil {
		return nil, err
	}

	return images, nil
}

// 根据用户名查询收藏的图片
func GetUserFavoritedImagesByUsername(db *gorm.DB, username string) ([]i.ImageInformation, error) {
	var images []i.ImageInformation
	err := db.Table("favoritedimage").Where("username = ?", username).Find(&images).Error // 使用 Find 而不是 First
	if err != nil {
		return nil, err
	}

	return images, nil
}

//-------------------------------------------------获取指定图像 或 （指定时间段内）所有图像

// 根据图片url查询相关图片
func GetImageByUrl(db *gorm.DB, url string) (*i.ImageInformation, error) {
	var image i.ImageInformation
	err := db.Table("imageinformation").Where("picture = ?", url).First(&image).Error
	if err != nil {
		return nil, err // 返回错误
	}

	return &image, nil // 返回指向image的指针
}

// 根据图片ID查询相关图片
func GetImageById(db *gorm.DB, id int) (*i.ImageInformation, error) {

	var image i.ImageInformation
	err := db.Table("imageinformation").Where("id = ?", id).First(&image).Error
	if err != nil {
		return nil, err // 返回错误
	}

	return &image, nil // 返回指向image的指针
}

// 根据图片的username查询图片
func GetImageByUsername(db *gorm.DB, username string) (*i.ImageInformation, error) {
	var image i.ImageInformation
	err := db.Table("imageinformation").Where("username = ?", username).First(&image).Error
	if err != nil {
		return nil, err // 返回错误
	}

	return &image, nil // 返回指向image的指针
}

// 查询指定时间段内的所有图像
func GetImagesInfoWithinTimeRange(db *gorm.DB, startTime, endTime time.Time) ([]i.ImageInformation, error) {
	var images []i.ImageInformation
	err := db.Table("imageinformation").
		Where("create_time BETWEEN ? AND ?", startTime, endTime).
		Find(&images).Error
	if err != nil {
		return nil, fmt.Errorf("查询图像列表时发生错误: %v", err)
	}
	return images, nil
}

// 获取所有图像信息并按id排序
func GetAllImagesInfo(db *gorm.DB) ([]i.ImageInformation, error) {
	var images []i.ImageInformation
	result := db.Table("imageinformation").Order("id ASC").Find(&images)
	if result.Error != nil {
		return nil, fmt.Errorf("查询图像列表时发生错误: %v", result.Error)
	}
	return images, nil
}
