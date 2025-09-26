package courses

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gateway/guards"
	"gateway/middlewares"
	"gateway/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CourseController struct {
	service *CourseService
}

func NewCourseController(s *CourseService) *CourseController {
	return &CourseController{
		service: s,
	}
}

func (c *CourseController) RegisterRoutes(r *gin.RouterGroup) {
	courseRoutes := r.Group("/courses")
	courseRoutes.Use(middlewares.AuthMiddleware())
	{
		courseRoutes.POST("/", middlewares.RequirePermission(guards.PermCourseCRUD), c.CreateCourse)
		courseRoutes.PUT("/:id", middlewares.RequirePermission(guards.PermCourseCRUD), c.UpdateCourseInfo)
		courseRoutes.PUT("/:id/publish", middlewares.RequirePermission(guards.PermCourseCRUD), c.PublishCourse)
		courseRoutes.GET("/:id", c.GetCourseByID)
	}
}

func (c *CourseController) CreateCourse(ctx *gin.Context) {
	var dto CreateCourseDTO
	userIdValue, _ := ctx.Get("userId")

	userId, ok := userIdValue.(uint)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid userId type"})
		return
	}

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(&utils.AppError{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	course, err := c.service.CreateCourse(dto, userId)
	if err != nil {
		ctx.Error(&utils.AppError{Status: http.StatusInternalServerError, Message: "Course creation failed"})
		return
	}

	ctx.JSON(http.StatusCreated, course)
}

func (c *CourseController) UpdateCourseInfo(ctx *gin.Context) {
	courseIdStr := ctx.Param("id")
	courseId, err := strconv.Atoi(courseIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course id"})
		return
	}

	userIdValue, _ := ctx.Get("userId")
	// roleValue, _ := ctx.Get("role")
	userId := userIdValue.(uint)
	// role := roleValue.(string)

	category := ctx.PostForm("category")
	level := ctx.PostForm("level")

	var bannerPath string
	bannerFile, err := ctx.FormFile("banner")
	if err == nil {
		bannerPath = fmt.Sprintf("uploads/banner/%d_%s", time.Now().Unix(), bannerFile.Filename)
		if err := ctx.SaveUploadedFile(bannerFile, bannerPath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save banner"})
			return
		}
	}

	var videoPath string
	videoFile, err := ctx.FormFile("video_preview")
	if err == nil {
		videoPath = fmt.Sprintf("uploads/videos/%d_%s", time.Now().Unix(), videoFile.Filename)
		if err := ctx.SaveUploadedFile(videoFile, videoPath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save video preview"})
			return
		}
	}

	updateData := map[string]interface{}{
		"banner":        bannerPath,
		"video_preview": videoPath,
		"category":      category,
		"level":         level,
	}

	// TODO: update role
	updatedCourse, err := c.service.UpdateCourseInfo(uint(courseId), updateData, userId, "admin")
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedCourse)
}

func (c *CourseController) PublishCourse(ctx *gin.Context) {
	courseIdStr := ctx.Param("id")
	courseId, err := strconv.Atoi(courseIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course id"})
		return
	}
	var body PublishCourseDTO
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course, err := c.service.TogglePublishCourse(uint(courseId), body.IsPublish)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Course status updated successfully",
		"data":    course,
	})
}

func (c *CourseController) GetCourseByID(ctx *gin.Context) {
	courseIdStr := ctx.Param("id")
	courseId, err := strconv.ParseUint(courseIdStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid course ID",
		})
		return
	}

	userId := ctx.MustGet("userId").(uint)
	role := ctx.GetString("role")
	course, err := c.service.GetCourseByID(uint(courseId), userId, role)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Course not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrive course",
		})
		return
	}
	ctx.JSON(http.StatusOK, course)
}
