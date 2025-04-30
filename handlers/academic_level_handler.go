package handlers

import (
	"github.com/gofiber/fiber/v2"
	"hello-fiber/repository"
	"hello-fiber/models"
)

type AcademicLevelHandler struct {
	repo *repository.AcademicLevelRepository
}

func NewAcademicLevelHandler(repo *repository.AcademicLevelRepository) *AcademicLevelHandler {
	return &AcademicLevelHandler{repo: repo}
}

func (h *AcademicLevelHandler) CreateLevel(c *fiber.Ctx) error {
	level := new(models.AcademicLevel)
	if err := c.BodyParser(level); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "خطأ في تحليل البيانات",
		})
	}

	if err := h.repo.Create(level); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في إنشاء المرحلة الأكاديمية",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(level)
}

func (h *AcademicLevelHandler) GetAllLevels(c *fiber.Ctx) error {
	levels, err := h.repo.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في استرجاع المراحل الأكاديمية",
		})
	}

	return c.JSON(levels)
}

func (h *AcademicLevelHandler) GetLevel(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف المرحلة غير صالح",
		})
	}

	level, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "المرحلة الأكاديمية غير موجودة",
		})
	}

	return c.JSON(level)
}

func (h *AcademicLevelHandler) UpdateLevel(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف المرحلة غير صالح",
		})
	}

	level := new(models.AcademicLevel)
	if err := c.BodyParser(level); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "خطأ في تحليل البيانات",
		})
	}

	level.ID = uint(id)
	if err := h.repo.Update(level); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في تحديث المرحلة الأكاديمية",
		})
	}

	return c.JSON(level)
}

func (h *AcademicLevelHandler) DeleteLevel(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف المرحلة غير صالح",
		})
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في حذف المرحلة الأكاديمية",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
