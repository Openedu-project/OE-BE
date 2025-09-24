package enrollments

import (
	"net/http"
	"strconv"

	"gateway/guards"
	"gateway/middlewares"

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
