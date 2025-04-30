package repository

import (
	"hello-fiber/models"

	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) Create(book *models.Book) error {
	return r.db.Create(book).Error
}

func (r *BookRepository) GetByID(id uint) (*models.Book, error) {
	book := new(models.Book)
	if err := r.db.Preload("Teacher").Preload("EducationPlatform").First(book, id).Error; err != nil {
		return nil, err
	}
	return book, nil
}

func (r *BookRepository) GetAll() ([]models.Book, error) {
	var books []models.Book
	if err := r.db.Preload("Teacher").Preload("EducationPlatform").Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (r *BookRepository) GetByTeacherID(teacherID uint) ([]models.Book, error) {
	var books []models.Book
	if err := r.db.Where("teacher_id = ?", teacherID).Preload("Teacher").Preload("EducationPlatform").Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (r *BookRepository) GetByPlatformID(platformID uint) ([]models.Book, error) {
	var books []models.Book
	if err := r.db.Where("platform_id = ?", platformID).Preload("Teacher").Preload("EducationPlatform").Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (r *BookRepository) Update(book *models.Book) error {
	return r.db.Save(book).Error
}

func (r *BookRepository) Delete(id uint) error {
	return r.db.Delete(&models.Book{}, id).Error
}
