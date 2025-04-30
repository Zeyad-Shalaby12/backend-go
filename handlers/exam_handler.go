package handlers

import (
	"hello-fiber/models"
	"hello-fiber/repository"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ExamHandler struct {
	repo *repository.ExamRepository
}

func NewExamHandler(repo *repository.ExamRepository) *ExamHandler {
	return &ExamHandler{repo: repo}
}

func (h *ExamHandler) CreateExam(c *fiber.Ctx) error {
	exam := new(models.Exam)
	if err := c.BodyParser(exam); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "خطأ في تحليل البيانات",
		})
	}

	if exam.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "اسم الامتحان مطلوب",
		})
	}

	if exam.TotalScore <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "الدرجة الكلية يجب أن تكون أكبر من صفر",
		})
	}

	if exam.Duration <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "مدة الامتحان يجب أن تكون أكبر من صفر",
		})
	}

	if err := h.repo.Create(exam); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في إنشاء الامتحان",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(exam)
}

func (h *ExamHandler) GetAllExams(c *fiber.Ctx) error {
	exams, err := h.repo.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في استرجاع الامتحانات",
		})
	}

	return c.JSON(exams)
}

func (h *ExamHandler) GetExam(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	exam, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "الامتحان غير موجود",
		})
	}

	return c.JSON(exam)
}

func (h *ExamHandler) GetExamsByCourse(c *fiber.Ctx) error {
	courseID, err := strconv.ParseUint(c.Params("courseId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف الدورة غير صالح",
		})
	}

	exams, err := h.repo.GetByCourse(uint(courseID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في استرجاع الامتحانات",
		})
	}

	return c.JSON(exams)
}

// GetQuestionBank يجلب جميع الأسئلة من الامتحانات التابعة لكورسات مرحلة دراسية محددة
func (h *ExamHandler) GetQuestionBank(c *fiber.Ctx) error {
	academicLevelID, err := strconv.ParseUint(c.Params("academicLevelId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف المرحلة الدراسية غير صالح",
		})
	}

	questions, err := h.repo.GetQuestionsByAcademicLevel(uint(academicLevelID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في استرجاع بنك الأسئلة",
		})
	}

	return c.JSON(fiber.Map{
		"questions": questions,
		"count":     len(questions),
	})
}

func (h *ExamHandler) UpdateExam(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	exam := new(models.Exam)
	if err := c.BodyParser(exam); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "خطأ في تحليل البيانات",
		})
	}

	exam.ID = uint(id)
	if err := h.repo.Update(exam); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في تحديث الامتحان",
		})
	}

	return c.JSON(exam)
}

func (h *ExamHandler) DeleteExam(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في حذف الامتحان",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
