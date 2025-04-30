package repository

import (
	"encoding/json"
	"hello-fiber/models"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

// GetDB يعيد قاعدة البيانات المستخدمة
func (r *PostRepository) GetDB() *gorm.DB {
	return r.db
}

// Create ينشئ منشور جديد
func (r *PostRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

// GetByID يجلب منشور بواسطة المعرف
func (r *PostRepository) GetByID(id uint) (*models.Post, error) {
	var post models.Post
	result := r.db.Preload("Student").Preload("Platform").First(&post, id)
	return &post, result.Error
}

// GetAll يجلب جميع المنشورات
func (r *PostRepository) GetAll() ([]models.Post, error) {
	var posts []models.Post
	result := r.db.Preload("Student").Preload("Platform").Order("created_at DESC").Find(&posts)
	return posts, result.Error
}

// GetByPlatform يجلب المنشورات الخاصة بمنصة معينة
func (r *PostRepository) GetByPlatform(platformID uint) ([]models.Post, error) {
	var posts []models.Post
	result := r.db.Preload("Student").Preload("Platform").Where("platform_id = ?", platformID).Order("created_at DESC").Find(&posts)
	return posts, result.Error
}

// GetByStudent يجلب المنشورات الخاصة بطالب معين
func (r *PostRepository) GetByStudent(studentID uint) ([]models.Post, error) {
	var posts []models.Post
	result := r.db.Preload("Student").Preload("Platform").Where("student_id = ?", studentID).Order("created_at DESC").Find(&posts)
	return posts, result.Error
}

// Update يحدث منشور
func (r *PostRepository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}

// Delete يحذف منشور
func (r *PostRepository) Delete(id uint) error {
	return r.db.Delete(&models.Post{}, id).Error
}

// IncrementLikes يزيد عدد الإعجابات بمقدار واحد
func (r *PostRepository) IncrementLikes(id uint) error {
	return r.db.Model(&models.Post{}).Where("id = ?", id).Update("likes_count", gorm.Expr("likes_count + ?", 1)).Error
}

// DecrementLikes ينقص عدد الإعجابات بمقدار واحد
func (r *PostRepository) DecrementLikes(id uint) error {
	return r.db.Model(&models.Post{}).Where("id = ? AND likes_count > 0", id).Update("likes_count", gorm.Expr("likes_count - ?", 1)).Error
}

// AddComment يضيف تعليق جديد إلى المنشور
func (r *PostRepository) AddComment(postID uint, studentID uint, studentName string, content string) error {
	// جلب المنشور
	post, err := r.GetByID(postID)
	if err != nil {
		return err
	}

	// إنشاء تعليق جديد
	newComment := models.Comment{
		StudentID: studentID,
		Name:      studentName,
		Content:   content,
		CreatedAt: time.Now().Unix(),
	}

	// تحويل التعليقات الحالية من JSON إلى مصفوفة
	var comments []models.Comment
	if post.Comments != nil {
		if err := json.Unmarshal(post.Comments, &comments); err != nil {
			return err
		}
	}

	// إضافة التعليق الجديد
	comments = append(comments, newComment)

	// تحويل المصفوفة إلى JSON
	commentsJSON, err := json.Marshal(comments)
	if err != nil {
		return err
	}

	// تحديث المنشور
	return r.db.Model(&models.Post{}).Where("id = ?", postID).Update("comments", datatypes.JSON(commentsJSON)).Error
}
