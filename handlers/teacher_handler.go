package handlers

import (
	"hello-fiber/models"
	"hello-fiber/repository"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type TeacherHandler struct {
	repo      *repository.TeacherRepository
	jwtSecret string
}

func NewTeacherHandler(repo *repository.TeacherRepository, jwtSecret string) *TeacherHandler {
	return &TeacherHandler{repo: repo, jwtSecret: jwtSecret}
}

func (h *TeacherHandler) Register(c *fiber.Ctx) error {
	teacher := new(models.Teacher)
	if err := c.BodyParser(teacher); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "خطأ في تحليل البيانات",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(teacher.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في معالجة كلمة المرور",
		})
	}
	teacher.Password = string(hashedPassword)

	if file, err := c.FormFile("image"); err == nil {
		ext := filepath.Ext(file.Filename)
		filename := "teacher_" + strconv.FormatUint(uint64(time.Now().Unix()), 10) + ext
		if err := c.SaveFile(file, "./uploads/"+filename); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "خطأ في حفظ الصورة",
			})
		}
		teacher.Image = filename
	}

	if err := h.repo.Create(teacher); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في إنشاء المعلم",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(teacher)
}

func (h *TeacherHandler) Login(c *fiber.Ctx) error {
	loginData := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "خطأ في تحليل البيانات",
		})
	}

	teacher, err := h.repo.GetByEmail(loginData.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "بيانات الدخول غير صحيحة",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(teacher.Password), []byte(loginData.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "بيانات الدخول غير صحيحة",
		})
	}

	claims := jwt.MapClaims{
		"id":    teacher.ID,
		"email": teacher.Email,
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
		"token": t,
	})
}

func (h *TeacherHandler) GetTeacher(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	teacher, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "المعلم غير موجود",
		})
	}

	return c.JSON(teacher)
}

func (h *TeacherHandler) UpdateTeacher(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	teacher := new(models.Teacher)
	if err := c.BodyParser(teacher); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "خطأ في تحليل البيانات",
		})
	}

	if file, err := c.FormFile("image"); err == nil {
		ext := filepath.Ext(file.Filename)
		filename := "teacher_" + strconv.FormatUint(uint64(time.Now().Unix()), 10) + ext
		if err := c.SaveFile(file, "./uploads/"+filename); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "خطأ في حفظ الصورة",
			})
		}
		teacher.Image = filename
	}

	teacher.ID = uint(id)
	if err := h.repo.Update(teacher); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في تحديث المعلم",
		})
	}

	return c.JSON(teacher)
}

func (h *TeacherHandler) DeleteTeacher(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في حذف المعلم",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
