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
	{
		enrollRoutesAuth.POST("/courses/:id/enroll", middlewares.RequirePermission(guards.PermEnrollInCourse), c.Enroll)
		enrollRoutesAuth.GET("/my-courses", middlewares.RequirePermission(guards.PermEnrollInCourse), c.GetMyCourses)
		enrollRoutesAuth.GET("/dashboard/learning-summary", middlewares.RequirePermission(guards.PermEnrollInCourse), c.GetDashboardSummary)
		enrollRoutesAuth.GET("/my-courses/:status", c.GetMyCoursesByStatus)
	}
}

func (c *EnrollmentController) Enroll(ctx *gin.Context) {
	// Get UserID from the authenticated context
	userIdValue, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: userId not found in context",
		})
		return
	}

	userId, ok := userIdValue.(uint)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid uesrId type in context",
		})
		return
	}

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
	userIdValue, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: userId not found in context",
		})
		return
	}
	userId, ok := userIdValue.(uint)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid userId type in context",
		})
		return
	}

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
	userIdValue, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: usreId not found in context",
		})
		return
	}

	userId, ok := userIdValue.(uint)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid userId tupe in context",
		})
		return
	}

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

	userIdValue, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorize: userId not found in context",
		})
		return
	}

	userId, ok := userIdValue.(uint)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid userId type in context",
		})
		return
	}

	switch statusStr {
	case "in-progress":
		page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
		courses, err := c.service.GetMyCoursesByStatus(userId, models.StatusInProgress, page, pageSize)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve in-progress courses",
			})
			return
		}
		ctx.JSON(http.StatusOK, courses)

	case "completed":
		page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
		courses, err := c.service.GetMyCoursesByStatus(userId, models.StatusCompleted, page, pageSize)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve completed courses",
			})
			return
		}
		ctx.JSON(http.StatusOK, courses)

	case "not-started":
		ctx.JSON(http.StatusOK, []CourseInfoDTO{})

	default:
		ctx.JSON(http.StatusBadRequest,
			gin.H{
				"error": "Invalid status value. Must be one of: in-progress,completed, not-started",
			})
	}
}
