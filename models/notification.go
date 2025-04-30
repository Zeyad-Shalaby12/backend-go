package models

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	Content             string            `json:"content" gorm:"not null"`
	EducationPlatformID uint              `json:"education_platform_id" gorm:"not null"`
	EducationPlatform   EducationPlatform `json:"education_platform" gorm:"foreignKey:EducationPlatformID"`
}


func MigrateNotification(db *gorm.DB) error {
	return db.AutoMigrate(&Notification{})
}