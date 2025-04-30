package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type EducationPlatform struct {
	gorm.Model
	Name              string         `json:"name" gorm:"not null"`
	PrimaryColor      string         `json:"primary_color"`
	IsUniversity      bool           `json:"is_university"`
	Status            string         `json:"status" gorm:"type:varchar(20);check:status IN ('active', 'not_active')"`
	HasOneTeacher     bool           `json:"has_one_teacher"`
	SupportDepartment bool           `json:"support_department"`
	ContactLinks      datatypes.JSON `json:"contact_links" gorm:"type:jsonb"`
	DarkMode          bool           `json:"dark_mode" gorm:"default:false"`
	GuestAllowed      bool           `json:"guest_allowed" gorm:"default:false"`
	MaxLoginCount     uint           `json:"max_login_count" gorm:"default:5"` // عدد مرات تسجيل الدخول المسموح بها
}

func MigrateEducationPlatforms(db *gorm.DB) error {
	return db.AutoMigrate(&EducationPlatform{})
}
