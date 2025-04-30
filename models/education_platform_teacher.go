package models

import (
	"errors"

	"gorm.io/gorm"
)

type EducationPlatformTeacher struct {
	gorm.Model
	TeacherID           uint              `json:"teacher_id" gorm:"not null"`
	EducationPlatformID uint              `json:"education_platform_id" gorm:"not null"`
	Teacher             Teacher           `json:"teacher" gorm:"foreignKey:TeacherID"`
	EducationPlatform   EducationPlatform `json:"education_platform" gorm:"foreignKey:EducationPlatformID"`
}

func MigratEducationPlatformTeacher(db *gorm.DB) error {
	return db.AutoMigrate(&EducationPlatformTeacher{})
}

func (ept *EducationPlatformTeacher) BeforeCreate(tx *gorm.DB) error {
	var platform EducationPlatform
	if err := tx.First(&platform, ept.EducationPlatformID).Error; err != nil {
		return err
	}

	if platform.HasOneTeacher {
		var count int64
		if err := tx.Model(&EducationPlatformTeacher{}).Where("education_platform_id = ?", ept.EducationPlatformID).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("هذه المنصة تقبل معلم واحد فقط")
		}
	}

	return nil
}
