package image_r

import (
	"errors"
	"fmt"
	"text-to-picture/models/image"
	"time"

	"gorm.io/gorm"
)

// // 根据用户ID查询相关图片
// func GetUserImagesByUserId(db *gorm.DB, userId int) ([]image.ImageInformation, error) {
// 	var images []image.ImageInformation
// 	err := db.Table("imageinformation").Where("userid = ?", userId).Find(&images).Error // 使用 Find 而不是 First
// 	if err != nil {
// 		return nil, err
// 	}

// 	return images, nil
// }

// // 根据用户ID查询收藏的图片
// func GetUserFavoritedImagesByUserId(db *gorm.DB, userId int) ([]image.ImageInformation, error) {
// 	var images []image.ImageInformation
// 	err := db.Table("favoritedimage").Where("userid = ?", userId).Find(&images).Error // 使用 Find 而不是 First
// 	if err != nil {
// 		return nil, err
// 	}

// 	return images, nil
// }

//-----------------------------------------------获取指定用户的图像

// 根据用户名查询相关图片
func GetUserImagesByUsername(db *gorm.DB, username string) ([]image.ImageInformation, error) {
	var images []image.ImageInformation
	err := db.Table("imageinformation").Where("username = ?", username).Find(&images).Error // 使用 Find 而不是 First
	if err != nil {
		return nil, err
	}

	return images, nil
}

// 根据用户名查询收藏的图片
func GetUserFavoritedImagesByUsername(db *gorm.DB, username string) ([]image.ImageInformation, error) {
	var images []image.ImageInformation
	err := db.Table("favoritedimage").Where("username = ?", username).Find(&images).Error // 使用 Find 而不是 First
	if err != nil {
		return nil, err
	}

	return images, nil
}

//-------------------------------------------------获取指定图像 或 （指定时间段内）所有图像

// 根据图片url查询相关图片
func GetImageByUrl(db *gorm.DB, url string) (*image.ImageInformation, error) {
	var image image.ImageInformation
	err := db.Table("imageinformation").Where("picture = ?", url).First(&image).Error
	if err != nil {
		return nil, err // 返回错误
	}

	return &image, nil // 返回指向image的指针
}

// 根据图片ID查询相关图片
func GetImageById(db *gorm.DB, id int) (*image.ImageInformation, error) {

	var image image.ImageInformation
	err := db.Table("imageinformation").Where("id = ?", id).First(&image).Error
	if err != nil {
		return nil, err // 返回错误
	}

	return &image, nil // 返回指向image的指针
}

// 根据图片的username查询图片
func GetImageByUsername(db *gorm.DB, username string) (*image.ImageInformation, error) {
	var image image.ImageInformation
	err := db.Table("imageinformation").Where("username = ?", username).First(&image).Error
	if err != nil {
		return nil, err // 返回错误
	}

	return &image, nil // 返回指向image的指针
}

// 根据图片ID查询相关图片
func GetFavoritedImageById(db *gorm.DB, id int) (*image.FavoritedImages, error) {

	var image image.FavoritedImages
	err := db.Table("favoritedimage").Where("id = ?", id).First(&image).Error
	if err != nil {
		return nil, err // 返回错误
	}

	return &image, nil // 返回指向image的指针
}

// 查询指定时间段内的所有图像
func GetImagesInfoWithinTimeRange(db *gorm.DB, username string, startTime, endTime time.Time) ([]image.ImageInformation, error) {
	var images []image.ImageInformation
	err := db.Table("imageinformation").
		Where("username = ?", username).
		Where("create_time BETWEEN ? AND ?", startTime, endTime).
		Find(&images).Error
	if err != nil {
		return nil, fmt.Errorf("查询图像列表时发生错误: %v", err)
	}
	return images, nil
}

// 获取所有图像信息并按id排序
func GetAllImagesInfo(db *gorm.DB) ([]image.ImageInformation, error) {
	var images []image.ImageInformation
	result := db.Table("imageinformation").Order("id ASC").Find(&images)
	if result.Error != nil {
		return nil, fmt.Errorf("查询图像列表时发生错误: %v", result.Error)
	}
	return images, nil
}

// 获取图像信息，并增加一个islike字段用于显示当前用户是否点赞过该图像
func GetAllImagesInfoWithLikeStatus(db *gorm.DB, username string) ([]image.ImageResponse, error) {
	var images []image.ImageInformation
	result := db.Table("imageinformation").Order("id ASC").Limit(100).Find(&images)
	if result.Error != nil {
		return nil, fmt.Errorf("查询图像列表时发生错误: %v", result.Error)
	}

	var num = 0
	imageResponses := make([]image.ImageResponse, len(images))
	for i, img := range images {
		imageResponses[i] = image.ImageResponse{
			ID:          img.ID,
			UserName:    img.UserName,
			Params:      img.Params,
			LikeCount:   img.LikeCount,
			Picture:     img.Picture,
			Create_time: img.Create_time,
			Isliked:     false, // 默认值为false
		}
		
		var imageLike image.ImageLike
		err := db.Where("username = ? AND picture = ?", username, img.Picture).First(&imageLike).Error
		if err == nil {
			imageResponses[i].Isliked = true
			num += 1
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果有其他类型的错误则返回
			return nil, fmt.Errorf("查询点赞记录时发生错误: %v", err)
		}
	}

	fmt.Println(num)

	return imageResponses, nil
}
