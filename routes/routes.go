package routes

import (
	"hello-fiber/handlers"
	"hello-fiber/middleware"
	"hello-fiber/repository"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, apiToken string) {
	api := app.Group("/api")
	
	// تطبيق middleware المصادقة على جميع مسارات API
	v1 := api.Group("/v1", middleware.AuthMiddleware(apiToken))

	// -------- المستودعات والمعالجات --------
	
	// 1. منصات التعليم
	educationPlatformRepo := repository.NewEducationPlatformRepository(db)
	educationPlatformHandler := handlers.NewEducationPlatformHandler(educationPlatformRepo)

	// 2. المدرسين
	teacherRepo := repository.NewTeacherRepository(db)
	teacherHandler := handlers.NewTeacherHandler(teacherRepo, "learnos")

	// 3. الطلاب
	studentRepo := repository.NewStudentRepository(db)
	studentHandler := handlers.NewStudentHandler(studentRepo, "learnos")

	// 4. علاقة المدرسين بالمنصات
	eptRepo := repository.NewEducationPlatformTeacherRepository(db)
	eptHandler := handlers.NewEducationPlatformTeacherHandler(eptRepo)

	// 5. المستويات الأكاديمية
	academicLevelRepo := repository.NewAcademicLevelRepository(db)
	academicLevelHandler := handlers.NewAcademicLevelHandler(academicLevelRepo)

	// 6. الكتب
	bookRepo := repository.NewBookRepository(db)
	bookHandler := handlers.NewBookHandler(bookRepo)

	// 7. الدورات
	courseRepo := repository.NewCourseRepository(db)
	courseHandler := handlers.NewCourseHandler(courseRepo)

	// 8. الاختبارات
	examRepo := repository.NewExamRepository(db)
	examHandler := handlers.NewExamHandler(examRepo)

	// 9. الفيديوهات
	videoRepo := repository.NewVideoRepository(db)
	videoHandler := handlers.NewVideoHandler(videoRepo)

	// 10. المنشورات
	postRepo := repository.NewPostRepository(db)
	postHandler := handlers.NewPostHandler(postRepo, studentRepo)

	// 11. الإشعارات
	notificationRepo := repository.NewNotificationRepository(db)
	notificationHandler := handlers.NewNotificationHandler(notificationRepo)

	// 12. تتبع تسجيل الدخول
	loginTrackingRepo := repository.NewLoginTrackingRepository(db)
	loginTrackingHandler := handlers.NewLoginTrackingHandler(loginTrackingRepo)

	// -------- مجموعات المسارات --------

	// 1. المنصات التعليمية
	platforms := v1.Group("/platforms")
	platforms.Post("/", educationPlatformHandler.CreatePlatform)
	platforms.Get("/", educationPlatformHandler.GetAllPlatforms)
	platforms.Get("/:id", educationPlatformHandler.GetPlatform)
	platforms.Put("/:id", educationPlatformHandler.UpdatePlatform)
	platforms.Delete("/:id", educationPlatformHandler.DeletePlatform)

	// 2. المستويات الأكاديمية
	levels := v1.Group("/academic-levels")
	levels.Post("/", academicLevelHandler.CreateLevel)
	levels.Get("/", academicLevelHandler.GetAllLevels)
	levels.Get("/:id", academicLevelHandler.GetLevel)
	levels.Put("/:id", academicLevelHandler.UpdateLevel)
	levels.Delete("/:id", academicLevelHandler.DeleteLevel)

	// 3. المستخدمين
	// 3.1 المدرسين
	teacherGroup := v1.Group("/teachers")
	teacherGroup.Post("/register", teacherHandler.Register)
	teacherGroup.Post("/login", teacherHandler.Login)
	teacherGroup.Get("/:id", teacherHandler.GetTeacher)
	teacherGroup.Put("/:id", teacherHandler.UpdateTeacher)
	teacherGroup.Delete("/:id", teacherHandler.DeleteTeacher)

	// 3.2 الطلاب
	students := v1.Group("/students")
	students.Post("/register", studentHandler.Register)
	students.Post("/login", studentHandler.Login)
	students.Get("/:id", studentHandler.GetStudent)
	students.Put("/:id", studentHandler.UpdateStudent)
	students.Delete("/:id", studentHandler.DeleteStudent)
	students.Post("/courses", studentHandler.AddCourse)
	students.Delete("/courses", studentHandler.RemoveCourse)
	students.Get("/:id/courses", studentHandler.GetStudentCourses)

	// 3.3 علاقة المدرسين بالمنصات
	eptGroup := v1.Group("/education-platform-teachers")
	eptGroup.Post("/", eptHandler.AssignTeacherToPlatform)
	eptGroup.Get("/platform/:platformId", eptHandler.GetPlatformTeachers)
	eptGroup.Get("/teacher/:teacherId", eptHandler.GetTeacherPlatforms)
	eptGroup.Delete("/platform/:platformId/teacher/:teacherId", eptHandler.RemoveTeacherFromPlatform)

	// 4. المحتوى التعليمي
	// 4.1 الكتب
	books := v1.Group("/books")
	books.Post("/", bookHandler.CreateBook)
	books.Get("/", bookHandler.GetAllBooks)
	books.Get("/:id", bookHandler.GetBook)
	books.Get("/teacher/:teacherId", bookHandler.GetBooksByTeacher)
	books.Get("/platform/:platformId", bookHandler.GetBooksByPlatform)
	books.Put("/:id", bookHandler.UpdateBook)
	books.Delete("/:id", bookHandler.DeleteBook)

	// 4.2 الدورات
	courses := v1.Group("/courses")
	courses.Post("/", courseHandler.CreateCourse)
	courses.Get("/", courseHandler.GetAllCourses)
	courses.Get("/:id", courseHandler.GetCourse)
	courses.Get("/academic-level/:academicLevelId", courseHandler.GetCoursesByAcademicLevel)
	courses.Put("/:id", courseHandler.UpdateCourse)
	courses.Delete("/:id", courseHandler.DeleteCourse)

	// 4.3 الاختبارات
	exams := v1.Group("/exams")
	exams.Post("/", examHandler.CreateExam)
	exams.Get("/", examHandler.GetAllExams)
	exams.Get("/:id", examHandler.GetExam)
	exams.Get("/course/:courseId", examHandler.GetExamsByCourse)
	exams.Get("/question-bank/:academicLevelId", examHandler.GetQuestionBank)
	exams.Put("/:id", examHandler.UpdateExam)
	exams.Delete("/:id", examHandler.DeleteExam)

	// 4.4 الفيديوهات
	videos := v1.Group("/videos")
	videos.Post("/", videoHandler.CreateVideo)
	videos.Get("/", videoHandler.GetAllVideos)
	videos.Get("/:id", videoHandler.GetVideo)
	videos.Get("/course/:courseId", videoHandler.GetVideosByCourse)
	videos.Put("/:id", videoHandler.UpdateVideo)
	videos.Delete("/:id", videoHandler.DeleteVideo)

	// 5. التفاعلات الاجتماعية
	// 5.1 المنشورات والتعليقات
	posts := v1.Group("/posts")
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

	// 5.2 الإشعارات
	notifications := v1.Group("/notifications")
	notifications.Post("/", notificationHandler.CreateNotification)
	notifications.Get("/", notificationHandler.GetAllNotifications)
	notifications.Get("/:id", notificationHandler.GetNotification)
	notifications.Get("/platform/:platformId", notificationHandler.GetNotificationsByPlatform)
	notifications.Put("/:id", notificationHandler.UpdateNotification)
	notifications.Delete("/:id", notificationHandler.DeleteNotification)

	// 6. تتبع النظام
	loginTracking := v1.Group("/login-tracking")
	loginTracking.Get("/platform/:platformId/limit", loginTrackingHandler.GetPlatformLoginLimit)
	loginTracking.Put("/platform/:platformId/limit", loginTrackingHandler.UpdatePlatformLoginLimit)
	loginTracking.Get("/student/:studentId/records", loginTrackingHandler.GetStudentLoginRecords)
	loginTracking.Get("/student/:studentId/platform/:platformId/count", loginTrackingHandler.GetStudentLoginCount)
}