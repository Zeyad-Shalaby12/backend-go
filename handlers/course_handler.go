package handlers

import (
	"hello-fiber/models"
	"hello-fiber/repository"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CourseHandler struct {
	repo *repository.CourseRepository
}

func NewCourseHandler(repo *repository.CourseRepository) *CourseHandler {
	return &CourseHandler{repo: repo}
}

func (h *CourseHandler) CreateCourse(c *fiber.Ctx) error {
	course := new(models.Course)
	if err := c.BodyParser(course); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "خطأ في تحليل البيانات",
		})
	}

	if course.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "اسم الدورة مطلوب",
		})
	}

	if course.Price < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "السعر يجب أن يكون أكبر من أو يساوي صفر",
		})
	}

	if err := h.repo.Create(course); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في إنشاء الدورة",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(course)
}

func (h *CourseHandler) GetAllCourses(c *fiber.Ctx) error {
	courses, err := h.repo.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في استرجاع الدورات",
		})
	}

	return c.JSON(courses)
}

func (h *CourseHandler) GetCourse(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	course, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "الدورة غير موجودة",
		})
	}

	return c.JSON(course)
}

func (h *CourseHandler) GetCoursesByAcademicLevel(c *fiber.Ctx) error {
	academicLevelID, err := strconv.ParseUint(c.Params("academicLevelId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف المرحلة الأكاديمية غير صالح",
		})
	}

	courses, err := h.repo.GetByAcademicLevel(uint(academicLevelID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في استرجاع الدورات",
		})
	}

	return c.JSON(courses)
}

func (h *CourseHandler) UpdateCourse(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	course := new(models.Course)
	if err := c.BodyParser(course); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "خطأ في تحليل البيانات",
		})
	}

	course.ID = uint(id)
	if err := h.repo.Update(course); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في تحديث الدورة",
		})
	}

	return c.JSON(course)
}

func (h *CourseHandler) DeleteCourse(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في حذف الدورة",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
