package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware يتحقق من وجود Bearer token صحيح في الـ Authorization header
func AuthMiddleware(staticToken string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// الحصول على قيمة Authorization header
		authHeader := c.Get("Authorization")

		// التحقق مما إذا كان الـ header موجود وبالتنسيق الصحيح
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "تم رفض الوصول: توكن المصادقة مفقود أو غير صالح",
			})
		}

		// استخراج التوكن من الـ header
		token := authHeader[7:]

		// التحقق من صحة التوكن
		if token != staticToken {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "تم رفض الوصول: توكن المصادقة غير صالح",
			})
		}

		// استمر إلى الـ handler التالي إذا كان التوكن صحيحاً
		return c.Next()
	}
}