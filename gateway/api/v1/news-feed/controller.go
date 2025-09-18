package newsfeed

import (
	"net/http"
	"strconv"

	"gateway/guards"
	"gateway/middlewares"

	"github.com/gin-gonic/gin"
)

type BlogController struct {
	service *BlogService
}

func NewBlogController(service *BlogService) *BlogController {
	return &BlogController{service: service}
}

func (c *BlogController) RegisterRoutes(r *gin.RouterGroup) {
	// public routes
	blogRoutes := r.Group("/blogs")
	{
		blogRoutes.GET("/")    // c.GetAllBlogs
		blogRoutes.GET("/:id") // c.GetBlogByID
	}

	// Auth routes
	blogRoutesAuth := r.Group("/blogs")
	blogRoutesAuth.Use(middlewares.AuthMiddleware())
	{
		blogRoutesAuth.POST("/", middlewares.RequirePermission(guards.Blog.Create), c.CreateBlog)
		// blogRoutesAuth.PUT("/:id", middlewares.RequirePermission(guards.Blog.Update), c.UpdateBlog)
		// blogRoutesAuth.DELETE("/:id", middlewares.RequirePermission(guards.Blog.Delete), c.DeleteBlog)
		blogRoutesAuth.POST("/:id/publish", middlewares.RequirePermission(guards.Blog.Update), c.RequestPublishBlog)
		blogRoutesAuth.POST("/:id/approve", middlewares.RequirePermission(guards.Blog.Update), c.ApproveBlog)
	}
}

func (c *BlogController) CreateBlog(ctx *gin.Context) {
	var req CreateBlogsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Lấy userID từ context (đã được middleware AuthMiddleware set vào)
	userIdValue, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not found in context",
		})
		return
	}

	authorID, ok := userIdValue.(uint)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user id type",
		})
		return
	}

	// Gọi service để tạo blog
	blog, err := c.service.CreateBlog(ctx, req, authorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, blog)
}

func (c *BlogController) RequestPublishBlog(ctx *gin.Context) {
	blogIDStr := ctx.Param("id")
	blogID, err := strconv.ParseUint(blogIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid blog ID",
		})
		return
	}

	role := ctx.GetString("role")

	blog, err := c.service.RequestPublish(ctx, uint(blogID), role)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, RequestPublishBlogResponse{
		Message: "Blog submitted for approval",
		Blog:    blog,
	})
}

func (c *BlogController) ApproveBlog(ctx *gin.Context) {
	blogIDStr := ctx.Param("id")
	blogID, err := strconv.ParseUint(blogIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid blog ID",
		})
		return
	}

	userIdValue, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found in context",
		})
		return
	}

	adminId, ok := userIdValue.(uint)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user id type",
		})
		return
	}

	var req ApproveBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	blog, err := c.service.ApproveBlog(ctx, uint(blogID), adminId, req)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, ApproveBlogResponse{
		Message: "Blog approval processed",
		Blog:    blog,
	})
}
