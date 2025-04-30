package handlers

import (
	"hello-fiber/models"
	"hello-fiber/repository"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type LoginTrackingHandler struct {
	repo *repository.LoginTrackingRepository
}

func NewLoginTrackingHandler(repo *repository.LoginTrackingRepository) *LoginTrackingHandler {
	return &LoginTrackingHandler{repo: repo}
}

// GetPlatformLoginLimit يحصل على حد تسجيل الدخول لمنصة معينة
func (h *LoginTrackingHandler) GetPlatformLoginLimit(c *fiber.Ctx) error {
	platformID, err := strconv.ParseUint(c.Params("platformId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف المنصة غير صالح",
		})
	}

	limit, err := h.repo.GetPlatformLoginLimit(uint(platformID))
	if err != nil {
		// إذا لم يتم العثور على حد، قم بإنشاء واحد افتراضي
		limit = &models.PlatformLoginLimit{
			EducationPlatformID: uint(platformID),
			MaxLoginCount:       5, // القيمة الافتراضية
		}
	}

	return c.JSON(limit)
}

// UpdatePlatformLoginLimit يحدث حد تسجيل الدخول لمنصة
func (h *LoginTrackingHandler) UpdatePlatformLoginLimit(c *fiber.Ctx) error {
	platformID, err := strconv.ParseUint(c.Params("platformId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف المنصة غير صالح",
		})
	}

	limit := new(models.PlatformLoginLimit)
	if err := c.BodyParser(limit); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "خطأ في تحليل البيانات",
		})
	}

	// التأكد من أن معرف المنصة في المسار يتطابق مع المعرف في الجسم
	limit.EducationPlatformID = uint(platformID)

	// التحقق من وجود حد مسبق
	existingLimit, err := h.repo.GetPlatformLoginLimit(uint(platformID))
	if err != nil {
		// إنشاء حد جديد
		if err := h.repo.CreatePlatformLoginLimit(limit); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "خطأ في إنشاء حد تسجيل الدخول",
			})
		}
	} else {
		// تحديث الحد الموجود
		existingLimit.MaxLoginCount = limit.MaxLoginCount
		if err := h.repo.UpdatePlatformLoginLimit(existingLimit); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "خطأ في تحديث حد تسجيل الدخول",
			})
		}
		limit = existingLimit
	}

	return c.JSON(limit)
}

// GetStudentLoginRecords يحصل على سجلات تسجيل دخول طالب معين
func (h *LoginTrackingHandler) GetStudentLoginRecords(c *fiber.Ctx) error {
	studentID, err := strconv.ParseUint(c.Params("studentId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف الطالب غير صالح",
		})
	}

	records, err := h.repo.GetStudentLoginRecords(uint(studentID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في استرجاع سجلات تسجيل الدخول",
		})
	}

	return c.JSON(records)
}

// GetStudentLoginCount يحصل على عدد مرات تسجيل دخول طالب لمنصة معينة
func (h *LoginTrackingHandler) GetStudentLoginCount(c *fiber.Ctx) error {
	studentID, err := strconv.ParseUint(c.Params("studentId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف الطالب غير صالح",
		})
	}

	platformID, err := strconv.ParseUint(c.Params("platformId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف المنصة غير صالح",
		})
	}

	count, err := h.repo.GetStudentLoginCount(uint(studentID), uint(platformID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في حساب عدد مرات تسجيل الدخول",
		})
	}

	// الحصول على الحد المسموح به
	limit, err := h.repo.GetPlatformLoginLimit(uint(platformID))
	if err != nil {
		// استخدام القيمة الافتراضية
		limit = &models.PlatformLoginLimit{
			MaxLoginCount: 5,
		}
	}

	return c.JSON(fiber.Map{
		"student_id":      studentID,
		"platform_id":     platformID,
		"login_count":     count,
		"max_login_count": limit.MaxLoginCount,
		"remaining":       int64(limit.MaxLoginCount) - count,
	})
}
