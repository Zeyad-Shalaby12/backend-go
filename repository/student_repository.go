package repository

import (
	"hello-fiber/models"

	"gorm.io/gorm"
)

type StudentRepository struct {
	db *gorm.DB
}

// GetDB يعيد قاعدة البيانات المستخدمة في المستودع
func (r *StudentRepository) GetDB() *gorm.DB {
	return r.db
}

func NewStudentRepository(db *gorm.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) Create(student *models.Student) error {
	return r.db.Create(student).Error
}

func (r *StudentRepository) GetAll() ([]models.Student, error) {
	var students []models.Student
	err := r.db.Preload("AcademicLevel").Preload("Platform").Preload("Courses").Find(&students).Error
	return students, err
}

func (r *StudentRepository) GetByID(id uint) (*models.Student, error) {
	var student models.Student
	err := r.db.Preload("AcademicLevel").Preload("Platform").Preload("Courses").First(&student, id).Error
	return &student, err
}

func (r *StudentRepository) Update(student *models.Student) error {
	return r.db.Save(student).Error
}

func (r *StudentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Student{}, id).Error
}

func (r *StudentRepository) AddCourse(studentID, courseID uint) error {
	return r.db.Exec("INSERT INTO student_courses (student_id, course_id) VALUES (?, ?)", studentID, courseID).Error
}

func (r *StudentRepository) RemoveCourse(studentID, courseID uint) error {
	return r.db.Exec("DELETE FROM student_courses WHERE student_id = ? AND course_id = ?", studentID, courseID).Error
}

func (r *StudentRepository) GetStudentCourses(studentID uint) ([]models.Course, error) {
	var courses []models.Course
	err := r.db.Joins("JOIN student_courses ON courses.id = student_courses.course_id").Where("student_courses.student_id = ?", studentID).Find(&courses).Error
	return courses, err
}

func (r *StudentRepository) GetByPhone(phone string) (*models.Student, error) {
	var student models.Student
	err := r.db.Preload("AcademicLevel").Preload("Platform").Preload("Courses").Where("phone = ?", phone).First(&student).Error
	return &student, err
}
