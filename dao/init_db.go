package dao

import "gorm.io/gorm"

func InitDB(db *gorm.DB) {
	db.AutoMigrate(&User{}, &Post{}, &Comment{})
}
