package repository

import (
	"hello-fiber/models"

	"gorm.io/gorm"
)

type VideoRepository struct {
	db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) *VideoRepository {
	return &VideoRepository{db: db}
}

func (r *VideoRepository) Create(video *models.Video) error {
	return r.db.Create(video).Error
}

func (r *VideoRepository) GetAll() ([]models.Video, error) {
	var videos []models.Video
	err := r.db.Preload("Course").Find(&videos).Error
	return videos, err
}

func (r *VideoRepository) GetByID(id uint) (*models.Video, error) {
	video := new(models.Video)
	err := r.db.Preload("Course").First(video, id).Error
	return video, err
}

func (r *VideoRepository) GetByCourse(courseID uint) ([]models.Video, error) {
	var videos []models.Video
	err := r.db.Preload("Course").Where("course_id = ?", courseID).Find(&videos).Error
	return videos, err
}

func (r *VideoRepository) Update(video *models.Video) error {
	return r.db.Save(video).Error
}

func (r *VideoRepository) Delete(id uint) error {
	return r.db.Delete(&models.Video{}, id).Error
}
