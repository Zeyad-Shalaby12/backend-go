package repository

import (
	"hello-fiber/models"

	"gorm.io/gorm"
)

type EducationPlatformTeacherRepository struct {
	db *gorm.DB
}

func NewEducationPlatformTeacherRepository(db *gorm.DB) *EducationPlatformTeacherRepository {
	return &EducationPlatformTeacherRepository{db: db}
}

func (r *EducationPlatformTeacherRepository) Create(ept *models.EducationPlatformTeacher) error {
	return r.db.Create(ept).Error
}

func (r *EducationPlatformTeacherRepository) GetByPlatformID(platformID uint) ([]models.EducationPlatformTeacher, error) {
	var epts []models.EducationPlatformTeacher
	if err := r.db.Where("education_platform_id = ?", platformID).Preload("Teacher").Find(&epts).Error; err != nil {
		return nil, err
	}
	return epts, nil
}

func (r *EducationPlatformTeacherRepository) GetByTeacherID(teacherID uint) ([]models.EducationPlatformTeacher, error) {
	var epts []models.EducationPlatformTeacher
	if err := r.db.Where("teacher_id = ?", teacherID).Preload("EducationPlatform").Find(&epts).Error; err != nil {
		return nil, err
	}
	return epts, nil
}

func (r *EducationPlatformTeacherRepository) Delete(platformID, teacherID uint) error {
	return r.db.Where("education_platform_id = ? AND teacher_id = ?", platformID, teacherID).Delete(&models.EducationPlatformTeacher{}).Error
}
