package model

import "github.com/jinzhu/gorm"

// User is a role, who has a login
type User struct {
	gorm.Model
	Login    string `json:"login"`
	Password string `json:"password"`
	Token    string `json:"token"`
}
