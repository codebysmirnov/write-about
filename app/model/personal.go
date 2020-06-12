package model

import "github.com/jinzhu/gorm"

// Personal is a personal data
type Personal struct {
	gorm.Model
	FirstName string `json:"firstName" gorm:"not null"`
	Lastname  string `json:"lastName"`
	Phone     string `json:"phone"`
	Email     string `json:"email" gorm:"unique;not null"`
	IDUser    uint   `json:"id_user" gorm:"index:id_user_personal_index"`
}
