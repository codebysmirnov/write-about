package model

import "github.com/jinzhu/gorm"

// DBMigrate - setup database tables
func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&User{}, &Personal{}, &Diary{})
	db.Model(&Personal{}).AddForeignKey("user_id", "users(id)", "NO ACTION", "NO ACTION")
	db.Model(&Diary{}).AddForeignKey("user_id", "users(id)", "NO ACTION", "NO ACTION")

	return db
}
