package models

type UserLogin struct {
	ID       int    `json:"id" gorm:"primarykey"`
	Email    string `json:"email" gorm:"unique;not null"`
	UserName string `json:"user_name" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
	Token    string `json:"token"`
}

type UserQuery struct {
	ID     int       `json:"id" gorm:"primarykey"`
	UserID int       `json:"user_id" gorm:"not null"`
	Params string    `json:"params"`
	Result string    `json:"result"`
	Time   string    `json:"time"`
	User   UserLogin `gorm:"foreignKey:UserID;references:ID"`
}
