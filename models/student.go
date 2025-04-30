package models

import (
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Name            string `json:"name" gorm:"not null"`
	Phone           string `json:"phone" gorm:"not null;unique"`
	ParentPhone     string `json:"parent_phone" gorm:"not null"`
	Gender          string `json:"gender" gorm:"not null"`     // ذكر أو أنثى
	StudyType       string `json:"study_type" gorm:"not null"` // منصة وسنتر أو منصة فقط
	Governorate     string `json:"governorate" gorm:"not null"`
	AcademicLevelID uint   `json:"academic_level_id" gorm:"not null"`
	PlatformID      uint   `json:"platform_id" gorm:"not null"`

	// العلاقات
	AcademicLevel AcademicLevel     `json:"academic_level" gorm:"foreignKey:AcademicLevelID"`
	Platform      EducationPlatform `json:"platform" gorm:"foreignKey:PlatformID"`
	Courses       []Course          `json:"courses" gorm:"many2many:student_courses;"`
}

func MigrateStudent(db *gorm.DB) error {
	return db.AutoMigrate(&Student{})
}