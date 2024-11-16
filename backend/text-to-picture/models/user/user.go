package user

import (
	"time"
)

type Login struct {
	ID          int       `json:"id" gorm:"primarykey"`
	Email       string    `json:"email" gorm:"unique;not null"`
	UserName    string    `json:"user_name" gorm:"not null"`
	Password    string    `json:"password" gorm:"not null"`
	Token       string    `json:"token"`
	Create_time time.Time `json:"create_time"`
}
type Register struct {
	ID          int       `json:"id" gorm:"primarykey"`
	Email       string    `json:"email" gorm:"unique;not null"`
	UserName    string    `json:"user_name" gorm:"not null"`
	Password    string    `json:"password" gorm:"not null"`
	Create_time time.Time `json:"create_time"`
}
