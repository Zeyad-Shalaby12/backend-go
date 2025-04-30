package models

import (
	"time"

	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	Name            string     `json:"name" gorm:"not null"`
	Description     string     `json:"description" gorm:"type:text"`
	Image           string     `json:"image"`
	Price           float64    `json:"price" gorm:"not null"`
	DiscountPrice   *float64   `json:"discount_price"`
	DiscountEndDate *time.Time `json:"discount_end_date"`
	AcademicLevelID uint       `json:"academic_level_id" gorm:"not null"`
	AcademicLevel   AcademicLevel
}
func MigratCourse(db *gorm.DB) error {
	return db.AutoMigrate(&Course{})
}
