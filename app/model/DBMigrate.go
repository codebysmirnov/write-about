package model

import "github.com/jinzhu/gorm"

// DBMigrate - setup database tables
func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&User{}, &Personal{}, &Diary{})
	db.Model(&Personal{}).AddForeignKey("id_user", "users(id)", "NO ACTION", "NO ACTION")
	db.Model(&Diary{}).AddForeignKey("id_user", "users(id)", "NO ACTION", "NO ACTION")

	return db
}
