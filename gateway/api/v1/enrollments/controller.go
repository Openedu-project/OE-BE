package enrollments

import (
	"net/http"
	"strconv"

	"gateway/guards"
	"gateway/middlewares"
	"gateway/models"

	"github.com/gin-gonic/gin"
)

type EnrollmentController struct {
	service *Service
}

func NewEnrollmentController(service *Service) *EnrollmentController {
	return &EnrollmentController{service: service}
}

func (c *EnrollmentController) RegisterRoutes(r *gin.RouterGroup) {
	enrollRoutesAuth := r.Group("/")
	enrollRoutesAuth.Use(middlewares.AuthMiddleware())
	enrollRoutesAuth.Use(middlewares.UserValidatorMiddleware())
	{
		enrollRoutesAuth.POST("/courses/:id/enroll", middlewares.RequirePermission(guards.PermEnrollInCourse), c.Enroll)
		enrollRoutesAuth.GET("/my-courses", middlewares.RequirePermission(guards.PermEnrollInCourse), c.GetMyCourses)
		enrollRoutesAuth.GET("/dashboard/learning-summary", middlewares.RequirePermission(guards.PermEnrollInCourse), c.GetDashboardSummary)
		enrollRoutesAuth.GET("/my-courses/:status", c.GetMyCoursesByStatus)
		enrollRoutesAuth.POST("/my-courses/:id/complete", middlewares.RequirePermission(guards.PermEnrollInCourse), c.CompleteCourse)
	}
}

func (c *EnrollmentController) Enroll(ctx *gin.Context) {
	userId := ctx.MustGet("userId").(uint)

	// Ger CourseId fromURL parameter
	courseIdStr := ctx.Param("id")
	courseId, err := strconv.ParseUint(courseIdStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid course ID",
		})
		return
	}

	// Call service to create the enrollment
	enrollment, err := c.service.CreateEnrollment(userId, uint(courseId))
	if err != nil {
		if err.Error() == "user already enrolled in this course" {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err.Error() == "cannot enroll in a non-existent course" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "An unexpected error occurred during enrollment",
		})
		return
	}

	response := EnrollmentResponseDTO{
		ID:       enrollment.ID,
		UserID:   enrollment.UserID,
		CourseID: enrollment.CourseID,
		Status:   string(enrollment.Status),
		CreateAt: enrollment.CreatedAt,
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *EnrollmentController) GetMyCourses(ctx *gin.Context) {
	userId := ctx.MustGet("userId").(uint)

	courses, err := c.service.GetMyCourses(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve courses",
		})
		return
	}

	ctx.JSON(http.StatusOK, courses)
}

func (c *EnrollmentController) GetDashboardSummary(ctx *gin.Context) {
	userId := ctx.MustGet("userId").(uint)

	summary, err := c.service.GetDashboardSummary(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retreve dashboard summary",
		})
		return
	}

	ctx.JSON(http.StatusOK, summary)
}

func (c *EnrollmentController) GetMyCoursesByStatus(ctx *gin.Context) {
	statusStr := ctx.Param("status")

	userId := ctx.MustGet("userId").(uint)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	var status models.UserCourseStatus

	switch statusStr {
	case "in-progress":
		status = models.StatusInProgress

	case "completed":
		status = models.StatusCompleted

	case "not-started":
		ctx.JSON(http.StatusOK, []CourseInfoDTO{})

	default:
		ctx.JSON(http.StatusBadRequest,
			gin.H{
				"error": "Invalid status value",
			})
		return
	}
	courses, err := c.service.GetMyCoursesByStatus(userId, status, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve courses by status",
		})
		return
	}
	ctx.JSON(http.StatusOK, courses)
}

func (c *EnrollmentController) CompleteCourse(ctx *gin.Context) {
	userId := ctx.MustGet("userId").(uint)

	courseIdStr := ctx.Param("id")
	courseId, err := strconv.ParseUint(courseIdStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid course ID",
		})
		return
	}

	_, err = c.service.CompletedCourse(userId, uint(courseId))
	if err != nil {
		if err.Error() == "user is not enrolled in this course" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err.Error() == "course is already completed" {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to complete course",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Course marked as completed successfully",
	})
}
