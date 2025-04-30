package handlers

import (
	"hello-fiber/models"
	"hello-fiber/repository"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PostHandler struct {
	repo        *repository.PostRepository
	studentRepo *repository.StudentRepository
}

func NewPostHandler(repo *repository.PostRepository, studentRepo *repository.StudentRepository) *PostHandler {
	return &PostHandler{repo: repo, studentRepo: studentRepo}
}

// CreatePost ينشئ منشور جديد
func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	post := new(models.Post)
	if err := c.BodyParser(post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "خطأ في تحليل البيانات",
		})
	}

	// التحقق من صحة البيانات
	if post.Content == "" || post.StudentID == 0 || post.PlatformID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "المحتوى ومعرف الطالب ومعرف المنصة مطلوبة",
		})
	}

	// التحقق من وجود الطالب
	_, err := h.studentRepo.GetByID(post.StudentID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "الطالب غير موجود",
		})
	}

	if err := h.repo.Create(post); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في إنشاء المنشور",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(post)
}

// GetPost يجلب منشور بواسطة المعرف
func (h *PostHandler) GetPost(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	post, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "المنشور غير موجود",
		})
	}

	return c.JSON(post)
}

// GetAllPosts يجلب جميع المنشورات
func (h *PostHandler) GetAllPosts(c *fiber.Ctx) error {
	posts, err := h.repo.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في جلب المنشورات",
		})
	}

	return c.JSON(posts)
}

// GetPostsByPlatform يجلب المنشورات الخاصة بمنصة معينة
func (h *PostHandler) GetPostsByPlatform(c *fiber.Ctx) error {
	platformID, err := strconv.ParseUint(c.Params("platformId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	posts, err := h.repo.GetByPlatform(uint(platformID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في جلب المنشورات",
		})
	}

	return c.JSON(posts)
}

// GetPostsByStudent يجلب المنشورات الخاصة بطالب معين
func (h *PostHandler) GetPostsByStudent(c *fiber.Ctx) error {
	studentID, err := strconv.ParseUint(c.Params("studentId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	posts, err := h.repo.GetByStudent(uint(studentID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في جلب المنشورات",
		})
	}

	return c.JSON(posts)
}

// UpdatePost يحدث منشور
func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	post := new(models.Post)
	if err := c.BodyParser(post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "خطأ في تحليل البيانات",
		})
	}

	post.ID = uint(id)
	if err := h.repo.Update(post); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في تحديث المنشور",
		})
	}

	return c.JSON(post)
}

// DeletePost يحذف منشور
func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في حذف المنشور",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// LikePost يزيد عدد الإعجابات بمقدار واحد
func (h *PostHandler) LikePost(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	if err := h.repo.IncrementLikes(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في زيادة الإعجابات",
		})
	}

	post, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في جلب المنشور",
		})
	}

	return c.JSON(fiber.Map{
		"likes_count": post.LikesCount,
	})
}

// UnlikePost ينقص عدد الإعجابات بمقدار واحد
func (h *PostHandler) UnlikePost(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف غير صالح",
		})
	}

	if err := h.repo.DecrementLikes(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في تنقيص الإعجابات",
		})
	}

	post, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في جلب المنشور",
		})
	}

	return c.JSON(fiber.Map{
		"likes_count": post.LikesCount,
	})
}

// AddComment يضيف تعليق جديد إلى المنشور
func (h *PostHandler) AddComment(c *fiber.Ctx) error {
	type CommentInput struct {
		PostID    uint   `json:"post_id"`
		StudentID uint   `json:"student_id"`
		Content   string `json:"content"`
	}

	input := new(CommentInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "خطأ في تحليل البيانات",
		})
	}

	// التحقق من صحة البيانات
	if input.PostID == 0 || input.StudentID == 0 || input.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "معرف المنشور ومعرف الطالب والمحتوى مطلوبة",
		})
	}

	// التحقق من وجود الطالب
	student, err := h.studentRepo.GetByID(input.StudentID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "الطالب غير موجود",
		})
	}

	// إضافة التعليق
	if err := h.repo.AddComment(input.PostID, input.StudentID, student.Name, input.Content); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في إضافة التعليق",
		})
	}

	// جلب المنشور بعد التحديث
	post, err := h.repo.GetByID(input.PostID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "خطأ في جلب المنشور",
		})
	}

	return c.JSON(post)
}
