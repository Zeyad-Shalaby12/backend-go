package handlers

import (
	"hello-fiber/models"
	"hello-fiber/repository"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type VideoHandler struct {
	repo *repository.VideoRepository
}

func NewVideoHandler(repo *repository.VideoRepository) *VideoHandler {
	return &VideoHandler{repo: repo}
}

func (h *VideoHandler) CreateVideo(c *fiber.Ctx) error {
	video := new(models.Video)
	if err := c.BodyParser(video); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "خطأ في تحليل البيانات",
		})
	}

	if video.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "اسم الفيديو مطلوب",
		})
	}

	if video.Type != "youtube" && video.Type != "server" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "نوع الفيديو يجب أن يكون youtube أو server",
		})
	}

	if video.URL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "رابط الفيديو مطلوب",
		})
	}

	if err := h.repo.Create(video); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في إنشاء الفيديو",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(video)
}

func (h *VideoHandler) GetAllVideos(c *fiber.Ctx) error {
	videos, err := h.repo.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في استرجاع الفيديوهات",
		})
	}

	return c.JSON(videos)
}

func (h *VideoHandler) GetVideo(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	video, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "الفيديو غير موجود",
		})
	}

	return c.JSON(video)
}

func (h *VideoHandler) GetVideosByCourse(c *fiber.Ctx) error {
	courseID, err := strconv.ParseUint(c.Params("courseId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف الدورة غير صالح",
		})
	}

	videos, err := h.repo.GetByCourse(uint(courseID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في استرجاع الفيديوهات",
		})
	}

	return c.JSON(videos)
}

func (h *VideoHandler) UpdateVideo(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	video := new(models.Video)
	if err := c.BodyParser(video); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "خطأ في تحليل البيانات",
		})
	}

	video.ID = uint(id)
	if err := h.repo.Update(video); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في تحديث الفيديو",
		})
	}

	return c.JSON(video)
}

func (h *VideoHandler) DeleteVideo(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في حذف الفيديو",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
