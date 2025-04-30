package models

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	Name     string `json:"name" gorm:"not null"`
	Type     string `json:"type" gorm:"not null"` // youtube or server
	URL      string `json:"url" gorm:"not null"`
	CourseID uint   `json:"course_id" gorm:"not null"`
	Course   Course `json:"course" gorm:"foreignKey:CourseID"`
}


func MigrateVideo(db *gorm.DB) error {
	return db.AutoMigrate(&Video{})
}