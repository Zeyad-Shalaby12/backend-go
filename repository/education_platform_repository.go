package repository

import (
	"hello-fiber/models"

	"gorm.io/gorm"
)

type EducationPlatformRepository struct {
	DB *gorm.DB
}

func NewEducationPlatformRepository(db *gorm.DB) *EducationPlatformRepository {
	return &EducationPlatformRepository{DB: db}
}

func (r *EducationPlatformRepository) Create(platform *models.EducationPlatform) error {
	return r.DB.Create(platform).Error
}

func (r *EducationPlatformRepository) GetAll() ([]models.EducationPlatform, error) {
	var platforms []models.EducationPlatform
	err := r.DB.Find(&platforms).Error
	return platforms, err
}

func (r *EducationPlatformRepository) GetByID(id uint) (*models.EducationPlatform, error) {
	var platform models.EducationPlatform
	err := r.DB.First(&platform, id).Error
	return &platform, err
}

func (r *EducationPlatformRepository) Update(platform *models.EducationPlatform) error {
	return r.DB.Save(platform).Error
}

func (r *EducationPlatformRepository) Delete(id uint) error {
	return r.DB.Delete(&models.EducationPlatform{}, id).Error
}
