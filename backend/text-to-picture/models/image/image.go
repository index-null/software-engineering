package models

import(
	u "gocode/backend/backend/text-to-picture/models/user"
)

type image struct {
	ID     int       `json:"id" gorm:"primarykey"`
	UserID string    `json:"user_id" gorm:"not null"`
	Result string    `json:"result"`
	User   u.UserLogin `gorm:"foreignKey:UserID;references:ID"`
}
