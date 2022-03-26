package models

type Auth struct {
	Id       int64  `json:"id" gorm:"primaryKey"`
	Email    string `json:"email" gorm:"index"`
	Password string `json:"password"`
}
