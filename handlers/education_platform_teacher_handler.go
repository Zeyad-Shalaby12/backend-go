package handlers

import (
	"hello-fiber/models"
	"hello-fiber/repository"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type EducationPlatformTeacherHandler struct {
	repo *repository.EducationPlatformTeacherRepository
}

func NewEducationPlatformTeacherHandler(repo *repository.EducationPlatformTeacherRepository) *EducationPlatformTeacherHandler {
	return &EducationPlatformTeacherHandler{repo: repo}
}

func (h *EducationPlatformTeacherHandler) AssignTeacherToPlatform(c *fiber.Ctx) error {
	ept := new(models.EducationPlatformTeacher)
	if err := c.BodyParser(ept); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "خطأ في تحليل البيانات",
		})
	}

	if err := h.repo.Create(ept); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(ept)
}

func (h *EducationPlatformTeacherHandler) GetPlatformTeachers(c *fiber.Ctx) error {
	platformID, err := strconv.ParseUint(c.Params("platformId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	epts, err := h.repo.GetByPlatformID(uint(platformID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في استرجاع المعلمين",
		})
	}

	return c.JSON(epts)
}

func (h *EducationPlatformTeacherHandler) GetTeacherPlatforms(c *fiber.Ctx) error {
	teacherID, err := strconv.ParseUint(c.Params("teacherId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	epts, err := h.repo.GetByTeacherID(uint(teacherID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في استرجاع المنصات",
		})
	}

	return c.JSON(epts)
}

func (h *EducationPlatformTeacherHandler) RemoveTeacherFromPlatform(c *fiber.Ctx) error {
	platformID, err := strconv.ParseUint(c.Params("platformId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف المنصة غير صالح",
		})
	}

	teacherID, err := strconv.ParseUint(c.Params("teacherId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف المعلم غير صالح",
		})
	}

	if err := h.repo.Delete(uint(platformID), uint(teacherID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في إزالة المعلم من المنصة",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
