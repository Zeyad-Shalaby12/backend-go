package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Comment يمثل بنية التعليق في المنشور
type Comment struct {
	StudentID uint   `json:"student_id"`
	Name      string `json:"name"`       // اسم الطالب
	Content   string `json:"content"`    // محتوى التعليق
	CreatedAt int64  `json:"created_at"` // وقت إنشاء التعليق
}

// Post يمثل منشور في النظام
type Post struct {
	gorm.Model
	Content    string         `json:"content" gorm:"not null"`
	LikesCount int            `json:"likes_count" gorm:"default:0"`
	Comments   datatypes.JSON `json:"comments" gorm:"type:jsonb;default:'[]'"`
	StudentID  uint           `json:"student_id" gorm:"not null"`
	PlatformID uint           `json:"platform_id" gorm:"not null"`

	// العلاقات
	Student  Student           `json:"student" gorm:"foreignKey:StudentID"`
	Platform EducationPlatform `json:"platform" gorm:"foreignKey:PlatformID"`
}

// MigratePost ينشئ جدول المنشورات في قاعدة البيانات
func MigratePost(db *gorm.DB) error {
	return db.AutoMigrate(&Post{})
}
