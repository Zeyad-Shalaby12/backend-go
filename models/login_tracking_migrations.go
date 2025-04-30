package models

import (
	"log"

	"gorm.io/gorm"
)

// MigrateLoginTracking تنفيذ الميجريشن لإنشاء جداول تتبع تسجيل الدخول
func MigrateLoginTracking(db *gorm.DB) error {
	log.Println("بدء تنفيذ ميجريشن جداول تتبع تسجيل الدخول...")

	// إنشاء جدول حدود تسجيل الدخول للمنصات
	if err := db.AutoMigrate(&PlatformLoginLimit{}); err != nil {
		log.Printf("خطأ في إنشاء جدول حدود تسجيل الدخول للمنصات: %v", err)
		return err
	}

	// إنشاء جدول سجلات تسجيل دخول الطلاب
	if err := db.AutoMigrate(&StudentLoginRecord{}); err != nil {
		log.Printf("خطأ في إنشاء جدول سجلات تسجيل دخول الطلاب: %v", err)
		return err
	}

	log.Println("تم تنفيذ ميجريشن جداول تتبع تسجيل الدخول بنجاح")
	return nil
}
