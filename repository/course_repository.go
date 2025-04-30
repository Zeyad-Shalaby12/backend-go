package repository

import (
	"hello-fiber/models"

	"gorm.io/gorm"
)

type CourseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) Create(course *models.Course) error {
	return r.db.Create(course).Error
}

func (r *CourseRepository) GetAll() ([]models.Course, error) {
	var courses []models.Course
	if err := r.db.Preload("AcademicLevel").Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *CourseRepository) GetByID(id uint) (*models.Course, error) {
	var course models.Course
	if err := r.db.Preload("AcademicLevel").First(&course, id).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *CourseRepository) GetByAcademicLevel(academicLevelID uint) ([]models.Course, error) {
	var courses []models.Course
	if err := r.db.Where("academic_level_id = ?", academicLevelID).Preload("AcademicLevel").Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *CourseRepository) Update(course *models.Course) error {
	return r.db.Save(course).Error
}

func (r *CourseRepository) Delete(id uint) error {
	return r.db.Delete(&models.Course{}, id).Error
}
