package repository

import (
	"gorm.io/gorm"
	"hello-fiber/models"
)

type AcademicLevelRepository struct {
	db *gorm.DB
}

func NewAcademicLevelRepository(db *gorm.DB) *AcademicLevelRepository {
	return &AcademicLevelRepository{db: db}
}

func (r *AcademicLevelRepository) Create(level *models.AcademicLevel) error {
	return r.db.Create(level).Error
}

func (r *AcademicLevelRepository) GetAll() ([]models.AcademicLevel, error) {
	var levels []models.AcademicLevel
	err := r.db.Preload("EducationPlatform").Find(&levels).Error
	return levels, err
}

func (r *AcademicLevelRepository) GetByID(id uint) (*models.AcademicLevel, error) {
	var level models.AcademicLevel
	err := r.db.Preload("EducationPlatform").First(&level, id).Error
	return &level, err
}

func (r *AcademicLevelRepository) Update(level *models.AcademicLevel) error {
	return r.db.Save(level).Error
}

func (r *AcademicLevelRepository) Delete(id uint) error {
	return r.db.Delete(&models.AcademicLevel{}, id).Error
}
