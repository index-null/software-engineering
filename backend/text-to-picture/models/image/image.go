package image

import (
	//u "text-to-picture/models/user"
	"time"
)

// type Image struct {
// 	ID          int       `json:"id" gorm:"primarykey"`
// 	UserID      string    `json:"user_id" gorm:"not null"`
// 	picture      string    `json:"picture"`
// 	Create_time time.Time `json:"create_time"`
// 	User        u.Login   `gorm:"foreignKey:UserID;references:ID"`
// }
// type QueryImage struct {
// 	ID          int       `json:"id" gorm:"primarykey"`
// 	picture      string    `json:"picture"`
// 	User        u.Login   `gorm:"foreignKey:UserID;references:ID"`
// 	Params      string    `json:"params"`
// 	Create_time time.Time `json:"create_time"`
// }

type ImageInformation struct {
	ID          int       `json:"id" gorm:"primarykey"`
	UserName    string    `json:"username" gorm:"column:username;not null"`
	Params      string    `json:"params"`
	LikeCount   int       `json:"likecount" gorm:"column:likecount"`
	Picture     string    `json:"picture"`
	Create_time time.Time `json:"create_time"`
	//User        u.UserInformation `gorm:"foreignKey:UserName;references:Username"`
}

func (ImageInformation) TableName() string {
	return "imageinformation"
}

type FavoritedImages struct {
	ID          int       `json:"id" gorm:"primarykey"`
	UserName    string    `json:"username" gorm:"column:username;not null"`
	Picture     string    `json:"picture"`
	Create_time time.Time `json:"create_time" gorm:"default:CURRENT_TIMESTAMP"`
}

func (FavoritedImages) TableName() string {
	return "favoritedimage"
}

type ImageLike struct {
	ID          int       `json:"id" gorm:"primarykey"`
	UserName    string    `json:"username" gorm:"column:username;not null"`
	Picture     string    `json:"picture"`
	Num         int       `json:"num"`
	Create_time time.Time `json:"create_time" gorm:"default:CURRENT_TIMESTAMP"`
}

func (ImageLike) TableName() string {
	return "imagelike"
}

type ImageResponse struct {
	ID          int       `json:"id"`
	UserName    string    `json:"username"`
	Params      string    `json:"params"`
	LikeCount   int       `json:"likecount"`
	Picture     string    `json:"picture"`
	Create_time time.Time `json:"create_time"`
	Isliked     bool      `json:"isliked"` // 注意这里修复了 JSON 标签
}

func (ImageResponse) TableName() string {
	return "imageresponse"
}
