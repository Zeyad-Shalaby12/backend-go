package handlers

import (
	"hello-fiber/models"
	"hello-fiber/repository"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type NotificationHandler struct {
	repo *repository.NotificationRepository
}

func NewNotificationHandler(repo *repository.NotificationRepository) *NotificationHandler {
	return &NotificationHandler{repo: repo}
}

func (h *NotificationHandler) CreateNotification(c *fiber.Ctx) error {
	notification := new(models.Notification)
	if err := c.BodyParser(notification); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "خطأ في تحليل البيانات",
		})
	}

	if notification.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "محتوى الإشعار مطلوب",
		})
	}

	if err := h.repo.Create(notification); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في إنشاء الإشعار",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(notification)
}

func (h *NotificationHandler) GetAllNotifications(c *fiber.Ctx) error {
	notifications, err := h.repo.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في استرجاع الإشعارات",
		})
	}

	return c.JSON(notifications)
}

func (h *NotificationHandler) GetNotification(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	notification, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "الإشعار غير موجود",
		})
	}

	return c.JSON(notification)
}

func (h *NotificationHandler) GetNotificationsByPlatform(c *fiber.Ctx) error {
	platformID, err := strconv.ParseUint(c.Params("platformId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف المنصة غير صالح",
		})
	}

	notifications, err := h.repo.GetByPlatform(uint(platformID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في استرجاع الإشعارات",
		})
	}

	return c.JSON(notifications)
}

func (h *NotificationHandler) UpdateNotification(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	notification := new(models.Notification)
	if err := c.BodyParser(notification); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "خطأ في تحليل البيانات",
		})
	}

	notification.ID = uint(id)
	if err := h.repo.Update(notification); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في تحديث الإشعار",
		})
	}

	return c.JSON(notification)
}

func (h *NotificationHandler) DeleteNotification(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في حذف الإشعار",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
