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
	Token       string    `json:"token"`
	Create_time time.Time `json:"create_time"`
}
gti
func (UserInformation) TableName() string {
	return "userinformation"
}

