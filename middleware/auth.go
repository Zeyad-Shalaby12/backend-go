package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

// AuthMiddleware يتحقق من وجود Bearer token صحيح في الـ Authorization header
func AuthMiddleware(staticToken string) fiber.Handler {
	// Print the expected token when the middleware is initialized (for debugging)
	log.Printf("AuthMiddleware initialized with token: %s", staticToken)
	
	return func(c *fiber.Ctx) error {
		// الحصول على قيمة Authorization header
		authHeader := c.Get("Authorization")
		
		// طباعة قيمة الهيدر للتصحيح
		log.Printf("Received Authorization header: %s", authHeader)

		// التحقق مما إذا كان الـ header موجود وبالتنسيق الصحيح
		if authHeader == "" {
			log.Println("Authorization header is missing")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "تم رفض الوصول: توكن المصادقة مفقود",
			})
		}

		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			log.Println("Authorization header does not have 'Bearer ' prefix")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "تم رفض الوصول: تنسيق التوكن غير صحيح (يجب أن يبدأ بـ 'Bearer ')",
			})
		}

		// استخراج التوكن من الـ header
		token := authHeader[7:]
		log.Printf("Extracted token: %s", token)
		log.Printf("Expected token: %s", staticToken)
		log.Printf("Token length: %d, Expected token length: %d", len(token), len(staticToken))

		// التحقق من صحة التوكن
		if token != staticToken {
			log.Println("Token validation failed - tokens do not match")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "تم رفض الوصول: توكن المصادقة غير صالح",
				"debug": fmt.Sprintf("Received: %s, Expected: %s", token, staticToken),
			})
		}

		log.Println("Token validation successful")
		// استمر إلى الـ handler التالي إذا كان التوكن صحيحاً
		return c.Next()
	}
}