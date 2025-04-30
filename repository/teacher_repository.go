package repository

import (
	"hello-fiber/models"

	"gorm.io/gorm"
)

type TeacherRepository struct {
	db *gorm.DB
}

func NewTeacherRepository(db *gorm.DB) *TeacherRepository {
	return &TeacherRepository{db: db}
}

func (r *TeacherRepository) Create(teacher *models.Teacher) error {
	return r.db.Create(teacher).Error
}

func (r *TeacherRepository) GetByID(id uint) (*models.Teacher, error) {
	teacher := new(models.Teacher)
	if err := r.db.First(teacher, id).Error; err != nil {
		return nil, err
	}
	return teacher, nil
}

func (r *TeacherRepository) GetByEmail(email string) (*models.Teacher, error) {
	teacher := new(models.Teacher)
	if err := r.db.Where("email = ?", email).First(teacher).Error; err != nil {
		return nil, err
	}
	return teacher, nil
}

func (r *TeacherRepository) GetAll() ([]models.Teacher, error) {
	var teachers []models.Teacher
	if err := r.db.Find(&teachers).Error; err != nil {
		return nil, err
	}
	return teachers, nil
}

func (r *TeacherRepository) Update(teacher *models.Teacher) error {
	return r.db.Save(teacher).Error
}

func (r *TeacherRepository) Delete(id uint) error {
	return r.db.Delete(&models.Teacher{}, id).Error
}
