package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// PlatformLoginLimit يحدد عدد مرات تسجيل الدخول المسموح بها لكل منصة
type PlatformLoginLimit struct {
	gorm.Model
	EducationPlatformID uint              `json:"education_platform_id" gorm:"not null"`
	MaxLoginCount       uint              `json:"max_login_count" gorm:"not null;default:5"` // عدد مرات تسجيل الدخول المسموح بها
	Platform            EducationPlatform `json:"platform" gorm:"foreignKey:EducationPlatformID"`
}

// StudentLoginRecord يسجل عمليات تسجيل دخول الطالب
type StudentLoginRecord struct {
	gorm.Model
	StudentID           uint      `json:"student_id" gorm:"not null"`
	EducationPlatformID uint      `json:"education_platform_id" gorm:"not null"`
	LoginTime           time.Time `json:"login_time" gorm:"not null"`
	IPAddress           string    `json:"ip_address"`
	UserAgent           string    `json:"user_agent"`

	// العلاقات
	Student  Student           `json:"student" gorm:"foreignKey:StudentID"`
	Platform EducationPlatform `json:"platform" gorm:"foreignKey:EducationPlatformID"`
}

// CheckLoginLimit يتحقق مما إذا كان الطالب قد تجاوز الحد المسموح به لتسجيل الدخول
func CheckLoginLimit(db *gorm.DB, studentID uint, platformID uint) error {
	// الحصول على حد تسجيل الدخول للمنصة
	var limit PlatformLoginLimit
	if err := db.Where("education_platform_id = ?", platformID).First(&limit).Error; err != nil {
		// إذا لم يتم العثور على حد، استخدم القيمة الافتراضية (5)
		limit.MaxLoginCount = 5
	}

	// حساب عدد مرات تسجيل الدخول للطالب على هذه المنصة
	var count int64
	if err := db.Model(&StudentLoginRecord{}).Where("student_id = ? AND education_platform_id = ?", studentID, platformID).Count(&count).Error; err != nil {
		return err
	}

	// التحقق من تجاوز الحد
	if uint(count) >= limit.MaxLoginCount {
		return errors.New("لقد تجاوزت الحد المسموح به لعدد مرات تسجيل الدخول لهذه المنصة")
	}

	return nil
}

// RecordLogin يسجل عملية تسجيل دخول جديدة للطالب
func RecordLogin(db *gorm.DB, studentID uint, platformID uint, ipAddress string, userAgent string) error {
	loginRecord := StudentLoginRecord{
		StudentID:           studentID,
		EducationPlatformID: platformID,
		LoginTime:           time.Now(),
		IPAddress:           ipAddress,
		UserAgent:           userAgent,
	}

	return db.Create(&loginRecord).Error
}
