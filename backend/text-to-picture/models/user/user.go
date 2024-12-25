package user

import (
	"time"
)

// UserInformation 定义用户信息的数据结构
type UserInformation struct {
	ID          int       `json:"id" gorm:"primarykey;autoIncrement"`
	Email       string    `json:"email" gorm:"unique;not null"`
	UserName    string    `json:"username" gorm:"column:username; not null"`
	Password    string    `json:"password" gorm:"not null"`
	Avatar_url  string    `json:"avatar_url"`
	Score       int       `json:"score"`
	Token       string    `json:"token"`
	Create_time time.Time `json:"create_time"`
}

// TableName 返回 UserInformation 数据表的名称
func (UserInformation) TableName() string {
	return "userinformation"
}

// UserScore 定义用户积分变动记录的数据结构
type UserScore struct {
	ID          int       `json:"id" gorm:"primarykey;autoIncrement"`
	Username    string    `json:"username"`
	Record      string    `json:"record"`
	Create_time time.Time `json:"create_time"`
}

// TableName 返回 UserScore 数据表的名称
func (UserScore) TableName() string {
	return "userscore"
}



