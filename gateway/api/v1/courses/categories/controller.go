package categories

import (
	"gateway/middlewares"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CourseCategoryController struct {
	service *CourseCategoryService
}

func NewCourseController(s *CourseCategoryService) *CourseCategoryController {
	return &CourseCategoryController{
		service: s,
	}
}

func (c *CourseCategoryController) RegisterRoutes(r *gin.RouterGroup) {
	courseCategoryRoutes := r.Group("/courses-categories")
	courseCategoryRoutes.Use(middlewares.AuthMiddleware())
	//TODO: add guards permission
	{
		courseCategoryRoutes.POST("/", c.CreateCategory)
		courseCategoryRoutes.GET("/", c.GetAllCategories)
		courseCategoryRoutes.GET("/:id", c.GetCategoryByID)
		courseCategoryRoutes.PUT("/:id", c.UpdateCategory)
		courseCategoryRoutes.DELETE("/:id", c.DeleteCategory)
	}
}

func NewCourseCategoryController(service *CourseCategoryService) *CourseCategoryController {
	return &CourseCategoryController{service: service}
}

func (c *CourseCategoryController) CreateCategory(ctx *gin.Context) {
	var dto CreateCourseCategoryDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := c.service.CreateCategory(dto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, category)
}

func (c *CourseCategoryController) GetCategoryByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	category, err := c.service.GetCategoryByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	ctx.JSON(http.StatusOK, category)
}

func (c *CourseCategoryController) GetAllCategories(ctx *gin.Context) {
	categories, err := c.service.GetAllCategories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, categories)
}

func (c *CourseCategoryController) UpdateCategory(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var dto UpdateCourseCategoryDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := c.service.UpdateCategory(uint(id), dto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, category)
}

func (c *CourseCategoryController) DeleteCategory(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := c.service.DeleteCategory(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}
