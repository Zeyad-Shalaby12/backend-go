package models

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name              string            `json:"name" gorm:"not null"`
	IsAvailable       bool              `json:"is_available" gorm:"default:true"`
	Image             string            `json:"image"`
	IsDigital         bool              `json:"is_digital"`
	Price             float64           `json:"price" gorm:"not null"`
	DiscountPrice     *float64          `json:"discount_price"`
	TeacherID         uint              `json:"teacher_id"`
	Teacher           Teacher           `json:"teacher" gorm:"foreignKey:TeacherID"`
	PlatformID        uint              `json:"platform_id"`
	EducationPlatform EducationPlatform `json:"education_platform" gorm:"foreignKey:PlatformID"`
}


func MigratBook(db *gorm.DB) error {
	return db.AutoMigrate(&Book{})
}
