package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

// Questions نوع مخصص للتعامل مع حقل الأسئلة JSONB
type Questions []map[string]interface{}

// Value تحويل Questions إلى قيمة يمكن تخزينها في قاعدة البيانات
func (q Questions) Value() (driver.Value, error) {
	return json.Marshal(q)
}

// Scan تحويل قيمة قاعدة البيانات إلى Questions
func (q *Questions) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("قيمة غير صالحة لحقل الأسئلة")
	}

	return json.Unmarshal(bytes, &q)
}

// Exam نموذج الامتحان
type Exam struct {
	gorm.Model
	Name       string    `json:"name" gorm:"not null"`
	TotalScore float64   `json:"total_score" gorm:"not null"`
	Duration   uint      `json:"duration" gorm:"not null"` // مدة الامتحان بالدقائق
	Questions  Questions `json:"questions" gorm:"type:jsonb"`
	CourseID   uint      `json:"course_id" gorm:"not null"`
	Course     Course    `json:"course" gorm:"foreignKey:CourseID"`
}


func MigrateExam(db *gorm.DB) error {
	return db.AutoMigrate(&Exam{})
}