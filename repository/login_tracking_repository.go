package repository

import (
	"hello-fiber/models"

	"gorm.io/gorm"
)

type LoginTrackingRepository struct {
	db *gorm.DB
}

func NewLoginTrackingRepository(db *gorm.DB) *LoginTrackingRepository {
	return &LoginTrackingRepository{db: db}
}

// GetPlatformLoginLimit يحصل على حد تسجيل الدخول لمنصة معينة
func (r *LoginTrackingRepository) GetPlatformLoginLimit(platformID uint) (*models.PlatformLoginLimit, error) {
	var limit models.PlatformLoginLimit
	err := r.db.Where("education_platform_id = ?", platformID).First(&limit).Error
	return &limit, err
}

// CreatePlatformLoginLimit ينشئ حد تسجيل دخول جديد لمنصة
func (r *LoginTrackingRepository) CreatePlatformLoginLimit(limit *models.PlatformLoginLimit) error {
	return r.db.Create(limit).Error
}

// UpdatePlatformLoginLimit يحدث حد تسجيل الدخول لمنصة
func (r *LoginTrackingRepository) UpdatePlatformLoginLimit(limit *models.PlatformLoginLimit) error {
	return r.db.Save(limit).Error
}

// GetStudentLoginRecords يحصل على سجلات تسجيل دخول طالب معين
func (r *LoginTrackingRepository) GetStudentLoginRecords(studentID uint) ([]models.StudentLoginRecord, error) {
	var records []models.StudentLoginRecord
	err := r.db.Where("student_id = ?", studentID).Order("login_time desc").Find(&records).Error
	return records, err
}

// GetStudentLoginCount يحصل على عدد مرات تسجيل دخول طالب لمنصة معينة
func (r *LoginTrackingRepository) GetStudentLoginCount(studentID uint, platformID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.StudentLoginRecord{}).Where("student_id = ? AND education_platform_id = ?", studentID, platformID).Count(&count).Error
	return count, err
}

// RecordStudentLogin يسجل عملية تسجيل دخول جديدة للطالب
func (r *LoginTrackingRepository) RecordStudentLogin(record *models.StudentLoginRecord) error {
	return r.db.Create(record).Error
}
