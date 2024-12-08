package user

import (
	"time"
)

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
type UserScore struct {
	ID          int       `json:"id" gorm:"primarykey;autoIncrement"`
	Username    string    `json:"username"`
	Record      string    `json:"record"`
	Create_time time.Time `json:"create_time"`
}

func (UserScore) TableName() string {
	return "userscore"
}
func (UserInformation) TableName() string {
	return "userinformation"
}
