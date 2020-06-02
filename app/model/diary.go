package model

import "github.com/jinzhu/gorm"

// Diary is a personal diary
type Diary struct {
	gorm.Model
	Year   int  `json:"year" gorm:"not null"`
	IDUser uint `json:"id_user" gorm:"index:id_user_index"`
}
