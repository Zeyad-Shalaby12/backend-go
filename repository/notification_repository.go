package repository

import (
	"hello-fiber/models"

	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

func (r *NotificationRepository) GetAll() ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.Preload("EducationPlatform").Find(&notifications).Error
	return notifications, err
}

func (r *NotificationRepository) GetByID(id uint) (*models.Notification, error) {
	notification := new(models.Notification)
	err := r.db.Preload("EducationPlatform").First(notification, id).Error
	return notification, err
}

func (r *NotificationRepository) GetByPlatform(platformID uint) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.Preload("EducationPlatform").Where("education_platform_id = ?", platformID).Find(&notifications).Error
	return notifications, err
}

func (r *NotificationRepository) Update(notification *models.Notification) error {
	return r.db.Save(notification).Error
}

func (r *NotificationRepository) Delete(id uint) error {
	return r.db.Delete(&models.Notification{}, id).Error
}
