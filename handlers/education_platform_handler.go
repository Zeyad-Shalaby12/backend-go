package handlers

import (
	"hello-fiber/models"
	"hello-fiber/repository"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type EducationPlatformHandler struct {
	repo *repository.EducationPlatformRepository
}

func NewEducationPlatformHandler(repo *repository.EducationPlatformRepository) *EducationPlatformHandler {
	return &EducationPlatformHandler{repo: repo}
}

func (h *EducationPlatformHandler) CreatePlatform(c *fiber.Ctx) error {
	platform := new(models.EducationPlatform)
	if err := c.BodyParser(platform); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "خطأ في تحليل البيانات",
		})
	}

	if err := h.repo.Create(platform); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في إنشاء المنصة",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(platform)
}

func (h *EducationPlatformHandler) GetAllPlatforms(c *fiber.Ctx) error {
	platforms, err := h.repo.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في استرجاع المنصات",
		})
	}

	return c.JSON(platforms)
}

func (h *EducationPlatformHandler) GetPlatform(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	platform, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "المنصة غير موجودة",
		})
	}

	return c.JSON(platform)
}

func (h *EducationPlatformHandler) UpdatePlatform(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	platform := new(models.EducationPlatform)
	if err := c.BodyParser(platform); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "خطأ في تحليل البيانات",
		})
	}

	platform.ID = uint(id)
	if err := h.repo.Update(platform); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في تحديث المنصة",
		})
	}

	return c.JSON(platform)
}

func (h *EducationPlatformHandler) DeletePlatform(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في حذف المنصة",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
