package user

import (
	"time"
)

type UserInformation struct {
	ID          int       `json:"id" gorm:"primarykey"`
	Email       string    `json:"email" gorm:"unique;not null"`
	UserName    string    `json:"user_name" gorm:"not null"`
	Password    string    `json:"password" gorm:"not null"`
	Avatar_nul  string    `json:"avatar_nul"`
	Token       string    `json:"token"`
	Create_time time.Time `json:"create_time"`
}

type Register struct {
	ID          int       `json:"id" gorm:"primarykey;autoIncrement"`
	Email       string    `json:"email" gorm:"unique;not null"`
	UserName    string    `json:"user_name" gorm:"not null"`
	Password    string    `json:"password" gorm:"not null"`
	Create_time time.Time `json:"create_time"`
}
