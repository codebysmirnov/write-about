package model

import "github.com/jinzhu/gorm"

// Diary is a personal diary
type Diary struct {
	gorm.Model
	Year   int  `json:"year" gorm:"not null"`
	UserID uint `json:"user_id" gorm:"index:user_id_index"`
}
