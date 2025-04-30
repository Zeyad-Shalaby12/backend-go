package handlers

import (
	"hello-fiber/models"
	"hello-fiber/repository"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type StudentHandler struct {
	repo      *repository.StudentRepository
	jwtSecret string
}

func NewStudentHandler(repo *repository.StudentRepository, jwtSecret string) *StudentHandler {
	return &StudentHandler{repo: repo, jwtSecret: jwtSecret}
}

func (h *StudentHandler) Register(c *fiber.Ctx) error {
	student := new(models.Student)
	if err := c.BodyParser(student); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "خطأ في تحليل البيانات",
		})
	}

	// التحقق من صحة البيانات
	if student.Name == "" || student.Phone == "" || student.ParentPhone == "" ||
		student.Gender == "" || student.StudyType == "" || student.Governorate == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "جميع الحقول مطلوبة",
		})
	}

	// التحقق من صحة الجنس
	if student.Gender != "ذكر" && student.Gender != "أنثى" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "الجنس يجب أن يكون ذكر أو أنثى",
		})
	}

	// التحقق من صحة نوع الدراسة
	if student.StudyType != "منصة وسنتر" && student.StudyType != "منصة فقط" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "نوع الدراسة يجب أن يكون منصة وسنتر أو منصة فقط",
		})
	}

	if err := h.repo.Create(student); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في إنشاء حساب الطالب",
		})
	}

	// إنشاء توكن JWT
	claims := jwt.MapClaims{
		"id":    student.ID,
		"phone": student.Phone,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في إنشاء التوكن",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"student": student,
		"token":   t,
	})
}

func (h *StudentHandler) Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Phone string `json:"phone"`
	}

	input := new(LoginInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "خطأ في تحليل البيانات",
		})
	}

	student, err := h.repo.GetByPhone(input.Phone)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "رقم الهاتف غير صحيح",
		})
	}

	// التحقق من عدد مرات تسجيل الدخول
	db := h.repo.GetDB()
	if err := models.CheckLoginLimit(db, student.ID, student.PlatformID); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// تسجيل عملية الدخول
	ipAddress := c.IP()
	userAgent := string(c.Request().Header.UserAgent())
	if err := models.RecordLogin(db, student.ID, student.PlatformID, ipAddress, userAgent); err != nil {
		// لا نريد منع تسجيل الدخول إذا فشل تسجيل العملية، لكن يجب تسجيل الخطأ
		// TODO: يمكن إضافة تسجيل الخطأ هنا
	}

	// إنشاء توكن JWT
	claims := jwt.MapClaims{
		"id":    student.ID,
		"phone": student.Phone,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في إنشاء التوكن",
		})
	}

	return c.JSON(fiber.Map{
		"student": student,
		"token":   t,
	})
}

func (h *StudentHandler) GetStudent(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	student, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "الطالب غير موجود",
		})
	}

	return c.JSON(student)
}

func (h *StudentHandler) UpdateStudent(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	student := new(models.Student)
	if err := c.BodyParser(student); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "خطأ في تحليل البيانات",
		})
	}

	student.ID = uint(id)
	if err := h.repo.Update(student); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في تحديث بيانات الطالب",
		})
	}

	return c.JSON(student)
}

func (h *StudentHandler) DeleteStudent(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في حذف الطالب",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *StudentHandler) AddCourse(c *fiber.Ctx) error {
	type CourseInput struct {
		StudentID uint `json:"student_id"`
		CourseID  uint `json:"course_id"`
	}

	input := new(CourseInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "خطأ في تحليل البيانات",
		})
	}

	if err := h.repo.AddCourse(input.StudentID, input.CourseID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في إضافة الدورة للطالب",
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *StudentHandler) RemoveCourse(c *fiber.Ctx) error {
	type CourseInput struct {
		StudentID uint `json:"student_id"`
		CourseID  uint `json:"course_id"`
	}

	input := new(CourseInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "خطأ في تحليل البيانات",
		})
	}

	if err := h.repo.RemoveCourse(input.StudentID, input.CourseID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في إزالة الدورة من الطالب",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *StudentHandler) GetStudentCourses(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	courses, err := h.repo.GetStudentCourses(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في استرجاع دورات الطالب",
		})
	}

	return c.JSON(courses)
}
