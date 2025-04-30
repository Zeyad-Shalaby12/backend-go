package routes

import (
	"hello-fiber/handlers"
	"hello-fiber/repository"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupRoutes يقوم بإعداد جميع مسارات التطبيق
// @title نظام إدارة المنصات التعليمية
// @version 1.0
// @description وثائق API لنظام إدارة المنصات التعليمية
// @contact.name فريق الدعم
// @contact.email support@eduplatform.com
// @host localhost:3000
// @BasePath /api
func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// إعداد API route prefix
	api := app.Group("/api")

	// إنشاء مستودع ومعالج المنصات التعليمية
	educationPlatformRepo := repository.NewEducationPlatformRepository(db)
	educationPlatformHandler := handlers.NewEducationPlatformHandler(educationPlatformRepo)

	// إنشاء مستودع ومعالج المعلمين
	teacherRepo := repository.NewTeacherRepository(db)
	teacherHandler := handlers.NewTeacherHandler(teacherRepo, "learnos")

	// إنشاء مستودع ومعالج العلاقات بين المعلمين والمنصات
	eptRepo := repository.NewEducationPlatformTeacherRepository(db)
	eptHandler := handlers.NewEducationPlatformTeacherHandler(eptRepo)

	// إعداد مسارات المنصات التعليمية
	platforms := api.Group("/platforms")
	// @Summary إنشاء منصة تعليمية جديدة
	// @Description إضافة منصة تعليمية جديدة إلى النظام
	// @Tags المنصات التعليمية
	// @Accept json
	// @Produce json
	// @Param platform body models.EducationPlatform true "بيانات المنصة"
	// @Success 201 {object} models.EducationPlatform
	// @Failure 400 {object} map[string]string
	// @Failure 500 {object} map[string]string
	// @Router /platforms [post]
	platforms.Post("/", educationPlatformHandler.CreatePlatform)
	// @Summary الحصول على جميع المنصات
	// @Description استرجاع قائمة بكافة المنصات التعليمية
	// @Tags المنصات التعليمية
	// @Produce json
	// @Success 200 {array} models.EducationPlatform
	// @Failure 500 {object} map[string]string
	// @Router /platforms [get]
	platforms.Get("/", educationPlatformHandler.GetAllPlatforms)
	// @Summary الحصول على منصة بواسطة ID
	// @Description استرجاع منصة تعليمية بواسطة المعرف
	// @Tags المنصات التعليمية
	// @Produce json
	// @Param id path int true "معرف المنصة"
	// @Success 200 {object} models.EducationPlatform
	// @Failure 400 {object} map[string]string
	// @Failure 404 {object} map[string]string
	// @Router /platforms/{id} [get]
	platforms.Get("/:id", educationPlatformHandler.GetPlatform)
	// @Summary تحديث منصة تعليمية
	// @Description تحديث بيانات منصة موجودة
	// @Tags المنصات التعليمية
	// @Accept json
	// @Produce json
	// @Param id path int true "معرف المنصة"
	// @Param platform body models.EducationPlatform true "بيانات التحديث"
	// @Success 200 {object} models.EducationPlatform
	// @Failure 400 {object} map[string]string
	// @Failure 404 {object} map[string]string
	// @Router /platforms/{id} [put]
	platforms.Put("/:id", educationPlatformHandler.UpdatePlatform)
	// @Summary حذف منصة تعليمية
	// @Description حذف منصة تعليمية بواسطة المعرف
	// @Tags المنصات التعليمية
	// @Param id path int true "معرف المنصة"
	// @Success 204
	// @Failure 400 {object} map[string]string
	// @Failure 404 {object} map[string]string
	// @Router /platforms/{id} [delete]
	platforms.Delete("/:id", educationPlatformHandler.DeletePlatform)

	// إعداد مسارات المعلمين
	teacherGroup := api.Group("/teachers")
	// @Summary تسجيل معلم جديد
	// @Description إضافة حساب معلم جديد إلى النظام
	// @Tags المعلمون
	// @Accept json
	// @Produce json
	// @Param teacher body models.Teacher true "بيانات التسجيل"
	// @Success 201 {object} models.Teacher
	// @Failure 400 {object} map[string]string
	// @Failure 500 {object} map[string]string
	// @Router /teachers/register [post]
	teacherGroup.Post("/register", teacherHandler.Register)
	// @Summary تسجيل دخول المعلم
	// @Description تسجيل دخول معلم موجود
	// @Tags المعلمون
	// @Accept json
	// @Produce json
	// @Param credentials body models.Login true "بيانات الدخول"
	// @Success 200 {object} map[string]string
	// @Failure 400 {object} map[string]string
	// @Failure 401 {object} map[string]string
	// @Router /teachers/login [post]
	teacherGroup.Post("/login", teacherHandler.Login)
	teacherGroup.Get("/:id", teacherHandler.GetTeacher)
	teacherGroup.Put("/:id", teacherHandler.UpdateTeacher)
	teacherGroup.Delete("/:id", teacherHandler.DeleteTeacher)

	// إعداد مسارات العلاقات بين المعلمين والمنصات
	eptGroup := api.Group("/education-platform-teachers")
	eptGroup.Post("/", eptHandler.AssignTeacherToPlatform)
	eptGroup.Get("/platform/:platformId", eptHandler.GetPlatformTeachers)
	eptGroup.Get("/teacher/:teacherId", eptHandler.GetTeacherPlatforms)
	// @Summary إزالة معلم من منصة
	// @Description إزالة ارتباط معلم بمنصة تعليمية
	// @Tags علاقات المنصات
	// @Param platformId path int true "معرف المنصة"
	// @Param teacherId path int true "معرف المعلم"
	// @Success 204
	// @Failure 400 {object} map[string]string
	// @Failure 500 {object} map[string]string
	// @Router /education-platform-teachers/platform/{platformId}/teacher/{teacherId} [delete]
	eptGroup.Delete("/platform/:platformId/teacher/:teacherId", eptHandler.RemoveTeacherFromPlatform)

	// إعداد مستودع ومعالج المراحل الأكاديمية
	academicLevelRepo := repository.NewAcademicLevelRepository(db)
	academicLevelHandler := handlers.NewAcademicLevelHandler(academicLevelRepo)

	// إعداد مسارات المراحل الأكاديمية
	levels := api.Group("/academic-levels")
	// @Summary إنشاء مرحلة أكاديمية
	// @Description إضافة مرحلة أكاديمية جديدة
	// @Tags المراحل الأكاديمية
	// @Accept json
	// @Produce json
	// @Param level body models.AcademicLevel true "بيانات المرحلة"
	// @Success 201 {object} models.AcademicLevel
	// @Failure 400 {object} map[string]string
	// @Failure 500 {object} map[string]string
	// @Router /academic-levels [post]
	levels.Post("/", academicLevelHandler.CreateLevel)
	levels.Get("/", academicLevelHandler.GetAllLevels)
	levels.Get("/:id", academicLevelHandler.GetLevel)
	levels.Put("/:id", academicLevelHandler.UpdateLevel)
	levels.Delete("/:id", academicLevelHandler.DeleteLevel)

	// إعداد مستودع ومعالج الكتب
	bookRepo := repository.NewBookRepository(db)
	bookHandler := handlers.NewBookHandler(bookRepo)

	// إعداد مسارات الكتب
	books := api.Group("/books")
	books.Post("/", bookHandler.CreateBook)
	books.Get("/", bookHandler.GetAllBooks)
	books.Get("/:id", bookHandler.GetBook)
	books.Get("/teacher/:teacherId", bookHandler.GetBooksByTeacher)
	books.Get("/platform/:platformId", bookHandler.GetBooksByPlatform)
	books.Put("/:id", bookHandler.UpdateBook)
	books.Delete("/:id", bookHandler.DeleteBook)

	// إعداد مستودع ومعالج الدورات
	courseRepo := repository.NewCourseRepository(db)
	courseHandler := handlers.NewCourseHandler(courseRepo)

	// إعداد مسارات الدورات
	courses := api.Group("/courses")
	courses.Post("/", courseHandler.CreateCourse)
	courses.Get("/", courseHandler.GetAllCourses)
	courses.Get("/:id", courseHandler.GetCourse)
	courses.Get("/academic-level/:academicLevelId", courseHandler.GetCoursesByAcademicLevel)
	courses.Put("/:id", courseHandler.UpdateCourse)
	courses.Delete("/:id", courseHandler.DeleteCourse)

	// إعداد مستودع ومعالج الامتحانات
	examRepo := repository.NewExamRepository(db)
	examHandler := handlers.NewExamHandler(examRepo)

	// إعداد مسارات الامتحانات
	exams := api.Group("/exams")
	exams.Post("/", examHandler.CreateExam)
	exams.Get("/", examHandler.GetAllExams)
	exams.Get("/:id", examHandler.GetExam)
	exams.Get("/course/:courseId", examHandler.GetExamsByCourse)
	exams.Get("/question-bank/:academicLevelId", examHandler.GetQuestionBank)
	exams.Put("/:id", examHandler.UpdateExam)
	exams.Delete("/:id", examHandler.DeleteExam)

	// إعداد مستودع ومعالج الفيديوهات
	videoRepo := repository.NewVideoRepository(db)
	videoHandler := handlers.NewVideoHandler(videoRepo)

	// إعداد مسارات الفيديوهات
	videos := api.Group("/videos")
	videos.Post("/", videoHandler.CreateVideo)

	// إنشاء مستودع ومعالج المنشورات
	postRepo := repository.NewPostRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	postHandler := handlers.NewPostHandler(postRepo, studentRepo)

	// إعداد مسارات المنشورات
	posts := api.Group("/posts")
	posts.Post("/", postHandler.CreatePost)
	posts.Get("/", postHandler.GetAllPosts)
	posts.Get("/:id", postHandler.GetPost)
	posts.Get("/platform/:platformId", postHandler.GetPostsByPlatform)
	posts.Get("/student/:studentId", postHandler.GetPostsByStudent)
	posts.Put("/:id", postHandler.UpdatePost)
	posts.Delete("/:id", postHandler.DeletePost)
	posts.Post("/:id/like", postHandler.LikePost)
	posts.Post("/:id/unlike", postHandler.UnlikePost)
	posts.Post("/comment", postHandler.AddComment)

	videos.Get("/", videoHandler.GetAllVideos)
	videos.Get("/:id", videoHandler.GetVideo)
	videos.Get("/course/:courseId", videoHandler.GetVideosByCourse)
	videos.Put("/:id", videoHandler.UpdateVideo)
	videos.Delete("/:id", videoHandler.DeleteVideo)

	// إعداد مستودع ومعالج الإشعارات
	notificationRepo := repository.NewNotificationRepository(db)
	notificationHandler := handlers.NewNotificationHandler(notificationRepo)

	// إعداد مسارات الإشعارات
	notifications := api.Group("/notifications")
	notifications.Post("/", notificationHandler.CreateNotification)
	notifications.Get("/", notificationHandler.GetAllNotifications)
	notifications.Get("/:id", notificationHandler.GetNotification)
	notifications.Get("/platform/:platformId", notificationHandler.GetNotificationsByPlatform)
	notifications.Put("/:id", notificationHandler.UpdateNotification)
	notifications.Delete("/:id", notificationHandler.DeleteNotification)

	// إعداد مستودع ومعالج الطلاب
	studentHandler := handlers.NewStudentHandler(studentRepo, "learnos")

	// إعداد مستودع ومعالج تتبع تسجيل الدخول
	loginTrackingRepo := repository.NewLoginTrackingRepository(db)
	loginTrackingHandler := handlers.NewLoginTrackingHandler(loginTrackingRepo)

	// إعداد مسارات الطلاب
	students := api.Group("/students")
	// @Summary تسجيل طالب جديد
	// @Description إضافة حساب طالب جديد إلى النظام
	// @Tags الطلاب
	// @Accept json
	// @Produce json
	// @Param student body models.Student true "بيانات التسجيل"
	// @Success 201 {object} models.Student
	// @Failure 400 {object} map[string]string
	// @Failure 500 {object} map[string]string
	// @Router /students/register [post]
	students.Post("/register", studentHandler.Register)
	// @Summary تسجيل دخول الطالب
	// @Description تسجيل دخول طالب موجود
	// @Tags الطلاب
	// @Accept json
	// @Produce json
	// @Param credentials body models.Login true "بيانات الدخول"
	// @Success 200 {object} map[string]string
	// @Failure 400 {object} map[string]string
	// @Failure 401 {object} map[string]string
	// @Router /students/login [post]
	students.Post("/login", studentHandler.Login)
	students.Get("/:id", studentHandler.GetStudent)
	students.Put("/:id", studentHandler.UpdateStudent)
	students.Delete("/:id", studentHandler.DeleteStudent)
	students.Post("/courses", studentHandler.AddCourse)
	students.Delete("/courses", studentHandler.RemoveCourse)
	students.Get("/:id/courses", studentHandler.GetStudentCourses)

	// إعداد مسارات تتبع تسجيل الدخول
	loginTracking := api.Group("/login-tracking")
	loginTracking.Get("/platform/:platformId/limit", loginTrackingHandler.GetPlatformLoginLimit)
	loginTracking.Put("/platform/:platformId/limit", loginTrackingHandler.UpdatePlatformLoginLimit)
	loginTracking.Get("/student/:studentId/records", loginTrackingHandler.GetStudentLoginRecords)
	loginTracking.Get("/student/:studentId/platform/:platformId/count", loginTrackingHandler.GetStudentLoginCount)
}
