package image

import (
	u "text-to-picture/models/user"
	"time"
)

type Image struct {
	ID     int     `json:"id" gorm:"primarykey"`
	UserID int  	`json:"user_id" gorm:"not null"`
	Result string  `json:"result"`
	User   u.Login `gorm:"foreignKey:UserID;references:ID"`
}
type QueryImage struct {
	ID     int       `json:"id" gorm:"primarykey"`
	Result string    `json:"result"`
	User   u.Login   `gorm:"foreignKey:UserID;references:ID"`
	Params string    `json:"params"`
	Time   time.Time `json:"time"`
}
