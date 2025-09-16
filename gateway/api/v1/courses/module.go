package courses

import (
	coursesCategory "gateway/api/v1/courses/categories"
	coursesSection "gateway/api/v1/courses/sections"
	"gateway/configs"
	"gateway/models"

	"github.com/gin-gonic/gin"
)

// Export instance
// var (
// 	CourseRepo *CourseRepository
// 	CourseSvc  *CourseService
// )

func InitModule(r *gin.Engine) {
	db := configs.DB
	if configs.Env.AppEnv != "production" {
		db.AutoMigrate(&models.Course{}, &models.CourseCategory{}, &models.CourseSection{}, &models.CourseLesson{})
	}

	CourseRepo := NewCourseRepository(db)
	CourseSvc := NewCourseService(CourseRepo)
	controller := NewCourseController(CourseSvc)

	api := r.Group("/api/v1")
	controller.RegisterRoutes(api)

	CourseCategoryRepo := coursesCategory.NewCourseCategoryRepository(db)
	CourseCategorySvc := coursesCategory.NewCourseCategoryService(CourseCategoryRepo)
	courseCategoryController := coursesCategory.NewCourseCategoryController(CourseCategorySvc)
	courseCategoryController.RegisterRoutes(api)

	CourseSectionRepo := coursesSection.NewCourseSectionRepository(db)
	CourseSectionSvc := coursesSection.NewCourseSectionService(CourseSectionRepo)
	courseSectionController := coursesSection.NewCourseSectionController(CourseSectionSvc)
	courseSectionController.RegisterRoutes(api)
}
