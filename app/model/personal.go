package model

import "github.com/jinzhu/gorm"

// Personal is a personal data
type Personal struct {
	gorm.Model
	FirstName string `json:"firstName" gorm:"not null"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Email     string `json:"email" gorm:"unique;not null"`
	UserID    uint   `json:"user_id" gorm:"index:user_id_personal_index"`
}
