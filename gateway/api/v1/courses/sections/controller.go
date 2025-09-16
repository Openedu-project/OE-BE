package sections

import (
	"gateway/middlewares"
	"gateway/models"
	"gateway/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CourseSectionController struct {
	service *CourseSectionService
}

func NewCourseSectionController(s *CourseSectionService) *CourseSectionController {
	return &CourseSectionController{
		service: s,
	}
}

func (c *CourseSectionController) RegisterRoutes(r *gin.RouterGroup) {
	courseRoutes := r.Group("/courses/:id/sections")
	courseRoutes.Use(middlewares.AuthMiddleware())
	{
		courseRoutes.POST("/", c.CreateCourseSection)
	}
}

func (c *CourseSectionController) CreateCourseSection(ctx *gin.Context) {
	var dto CreateCourseSectionDTO
	id := ctx.Param("id")
	courseId, _ := strconv.Atoi(id)

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(&utils.AppError{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	course, err := c.service.CreateCourseSection(dto, uint(courseId))
	if err != nil {
		ctx.Error(&utils.AppError{Status: http.StatusInternalServerError, Message: "CourseSection creation failed"})
		return
	}

	ctx.JSON(http.StatusCreated, course)
}

func (c *CourseSectionController) GetByID(ctx *gin.Context) (*models.CourseSection, error) {
	id := ctx.Param("id")
	sectionId, _ := strconv.Atoi(id)
	section, err := c.service.GetCourseSectionByID(uint(sectionId))

	if err != nil {
		return nil, err
	}

	return section, nil
}
