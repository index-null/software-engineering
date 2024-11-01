package user

type Login struct {
	ID       int    `json:"id" gorm:"primarykey"`
	Email    string `json:"email" gorm:"unique;not null"`
	UserName string `json:"user_name" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
	Token    string `json:"token"`
}
type Register struct {
	ID       int    `json:"id" gorm:"primarykey"`
	Email    string `json:"email" gorm:"unique;not null"`
	UserName string `json:"user_name" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
}
