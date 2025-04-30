package repository

import (
	"hello-fiber/models"

	"gorm.io/gorm"
)

type ExamRepository struct {
	db *gorm.DB
}

func NewExamRepository(db *gorm.DB) *ExamRepository {
	return &ExamRepository{db: db}
}

func (r *ExamRepository) Create(exam *models.Exam) error {
	return r.db.Create(exam).Error
}

func (r *ExamRepository) GetAll() ([]models.Exam, error) {
	var exams []models.Exam
	err := r.db.Preload("Course").Find(&exams).Error
	return exams, err
}

func (r *ExamRepository) GetByID(id uint) (*models.Exam, error) {
	var exam models.Exam
	err := r.db.Preload("Course").First(&exam, id).Error
	return &exam, err
}

func (r *ExamRepository) GetByCourse(courseID uint) ([]models.Exam, error) {
	var exams []models.Exam
	err := r.db.Preload("Course").Where("course_id = ?", courseID).Find(&exams).Error
	return exams, err
}

func (r *ExamRepository) Update(exam *models.Exam) error {
	return r.db.Save(exam).Error
}

func (r *ExamRepository) Delete(id uint) error {
	return r.db.Delete(&models.Exam{}, id).Error
}

// GetQuestionsByAcademicLevel يجلب جميع الأسئلة من الامتحانات التابعة لكورسات مرحلة دراسية محددة
func (r *ExamRepository) GetQuestionsByAcademicLevel(academicLevelID uint) ([]map[string]interface{}, error) {
	var exams []models.Exam

	// جلب جميع الامتحانات المرتبطة بالكورسات التابعة للمرحلة الدراسية المحددة
	err := r.db.Joins("JOIN courses ON exams.course_id = courses.id").Where("courses.academic_level_id = ?", academicLevelID).Preload("Course").Preload("Course.AcademicLevel").Find(&exams).Error
	if err != nil {
		return nil, err
	}

	// تجميع الأسئلة من جميع الامتحانات
	result := make([]map[string]interface{}, 0)
	for _, exam := range exams {
		for _, question := range exam.Questions {
			// إضافة معلومات الكورس والامتحان إلى كل سؤال
			question["exam_name"] = exam.Name
			question["course_name"] = exam.Course.Name
			question["academic_level_name"] = exam.Course.AcademicLevel.Name
			result = append(result, question)
		}
	}

	return result, nil
}
