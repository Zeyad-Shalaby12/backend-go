package models

import (
	"gorm.io/gorm"
)

type AcademicLevel struct {
	gorm.Model
	Name                string            `json:"name" gorm:"not null"`
	Description         string            `json:"description"`
	Image               string            `json:"image"`
	EducationPlatformID uint              `json:"education_platform_id"`
	EducationPlatform   EducationPlatform `json:"education_platform" gorm:"foreignKey:EducationPlatformID"`
}

func MigrateAcademicLevels(db *gorm.DB) error {
	return db.AutoMigrate(&AcademicLevel{})
}
