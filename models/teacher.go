package models

import (
	"gorm.io/gorm"
)

type Teacher struct {
	gorm.Model
	Name     string `json:"name" gorm:"not null"`
	Image    string `json:"image"`
	Phone    string `json:"phone" gorm:"unique"`
	Password string `json:"password" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
}

func MigrateTeacher(db *gorm.DB) error {
	return db.AutoMigrate(&Teacher{})
}
