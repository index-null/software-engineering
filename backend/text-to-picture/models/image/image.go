package image

import (
	//u "text-to-picture/models/user"
	"time"
)

// ImageInformation 定义图像信息的数据结构
type ImageInformation struct {
	ID          int       `json:"id" gorm:"primarykey"`
	UserName    string    `json:"username" gorm:"column:username;not null"`
	Params      string    `json:"params"`
	LikeCount   int       `json:"likecount" gorm:"column:likecount"`
	Picture     string    `json:"picture"`
	Create_time time.Time `json:"create_time"`
}

// TableName 返回数据库中表的名称
func (ImageInformation) TableName() string {
	return "imageinformation"
}

// FavoritedImages 定义用户收藏的图像信息的数据结构
type FavoritedImages struct {
	ID          int       `json:"id" gorm:"primarykey"`
	UserName    string    `json:"username" gorm:"column:username;not null"`
	Picture     string    `json:"picture"`
	Create_time time.Time `json:"create_time" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName 返回数据库中表的名称
func (FavoritedImages) TableName() string {
	return "favoritedimage"
}

// ImageLike 定义用户对图像的点赞信息的数据结构
type ImageLike struct {
	ID          int       `json:"id" gorm:"primarykey"`
	UserName    string    `json:"username" gorm:"column:username;not null"`
	Picture     string    `json:"picture"`
	Num         int       `json:"num"`
	Create_time time.Time `json:"create_time" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName 返回数据库中表的名称
func (ImageLike) TableName() string {
	return "imagelike"
}

// ImageResponse 定义图像响应的数据结构，包含用户的点赞状态
type ImageResponse struct {
	ID          int       `json:"id"`
	UserName    string    `json:"username"`
	Params      string    `json:"params"`
	LikeCount   int       `json:"likecount"`
	Picture     string    `json:"picture"`
	Create_time time.Time `json:"create_time"`
	Isliked     bool      `json:"isliked"` 
}

// TableName 返回数据库中表的名称
func (ImageResponse) TableName() string {
	return "imageresponse"
}
