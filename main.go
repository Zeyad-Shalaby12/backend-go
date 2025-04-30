package main

import (
	"hello-fiber/database"
	"hello-fiber/routes"
	"hello-fiber/storage"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// تحديد التوكن الثابت من متغير البيئة أو استخدام القيمة الافتراضية

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// تنفيذ جميع الميجريشنز
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "حدث خطأ في الخادم: " + err.Error(),
			})
		},
	})

	// إضافة middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// إعداد جميع المسارات مع تمرير توكن المصادقة
	routes.SetupRoutes(app, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("بدء تشغيل الخادم على المنفذ %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("فشل في بدء الخادم: %v", err)
	}
}