package models

type image struct {
	ID     int       `json:"id" gorm:"primarykey"`
	UserID string    `json:"user_id" gorm:"not null"`
	Result string    `json:"result"`
	User   UserLogin `gorm:"foreignKey:UserID;references:ID"`
}
