package main

import (
	"hello-fiber/database"
	"hello-fiber/routes"
	"hello-fiber/storage"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

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

	app := fiber.New()

	// إعداد جميع المسارات
	routes.SetupRoutes(app, db)

	app.Listen(":8000")
}