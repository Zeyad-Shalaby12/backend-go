package database

import (
	"hello-fiber/models"
	"log"

	"gorm.io/gorm"
)

// RunMigrations تنفيذ جميع الميجريشنز الخاصة بالتطبيق
func RunMigrations(db *gorm.DB) error {
	log.Println("بدء تنفيذ الميجريشنز...")

	// تنفيذ الميجريشن لإنشاء الجداول
	if err := models.MigrateEducationPlatforms(db); err != nil {
		return err
	}

	if err := models.MigrateTeacher(db); err != nil {
		return err
	}

	if err := models.MigratEducationPlatformTeacher(db); err != nil {
		return err
	}

	if err := models.MigrateAcademicLevels(db); err != nil {
		return err
	}

	if err := models.MigratBook(db); err != nil {
		return err
	}

	if err := models.MigratCourse(db); err != nil {
		return err
	}

	if err := models.MigrateExam(db); err != nil {
		return err
	}

	if err := models.MigrateNotification(db); err != nil {
		return err
	}

	if err := models.MigrateVideo(db); err != nil {
		return err
	}

	if err := models.MigrateStudent(db); err != nil {
		return err
	}

	if err := models.MigrateLoginTracking(db); err != nil {
		return err
	}

	if err := models.MigratePost(db); err != nil {
		return err
	}

	log.Println("تم تنفيذ جميع الميجريشنز بنجاح")
	return nil
}
