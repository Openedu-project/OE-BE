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
		blogRoutes.GET("/", c.ListBlogs)        // list published blogs
		blogRoutes.GET("/:id", c.GetBlogDetail) // blog detail (published)
	}

	// Auth routes (must be behind AuthMiddleware)
	blogRoutesAuth := r.Group("/blogs")
	blogRoutesAuth.Use(middlewares.AuthMiddleware())
	{
		blogRoutesAuth.POST("/", middlewares.RequirePermission(guards.BlogCreate), c.CreateBlog)
		blogRoutesAuth.PUT("/:id", middlewares.RequirePermission(guards.BlogUpdate), c.UpdateBlog)
		blogRoutesAuth.DELETE("/:id", middlewares.RequirePermission(guards.BlogDelete), c.DeleteBlog)
		blogRoutesAuth.POST("/:id/publish", middlewares.RequirePermission(guards.BlogPublish), c.RequestPublishBlog)
		blogRoutesAuth.POST("/:id/approve", middlewares.RequirePermission(guards.BlogPublish), c.ApproveBlog)

		// My blogs
		blogRoutesAuth.GET("/my", c.MyBlogs)
	}

	// Category management (admin only) - require PermUserManage (admin/sysadmin)
	catGroup := blogRoutesAuth.Group("/categories")
	{
		catGroup.POST("/", middlewares.RequirePermission(guards.PermUserManage), c.CreateCategory)
		catGroup.PUT("/:id", middlewares.RequirePermission(guards.PermUserManage), c.UpdateCategory)
		catGroup.DELETE("/:id", middlewares.RequirePermission(guards.PermUserManage), c.DeleteCategory)
	}

	// Public categories list
	catPublic := r.Group("/categories")
	{
		catPublic.GET("/", c.ListCategory)
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

	userIdValue, _ := ctx.Get("userId")
	authorID := userIdValue.(uint)

	blog, err := c.service.CreateBlog(ctx, req, authorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if ctx.Query("publish") == "true" {
		blog, err = c.service.RequestPublish(ctx, blog.ID, authorID)
		if err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message": "Blog created and sumitted for approval",
			"blog":    blog,
		})
		return
	}

	// Chỉ lưu Draft
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Blog save as draft",
		"blog":    blog,
	})
}

func (c *BlogController) RequestPublishBlog(ctx *gin.Context) {
	blogID, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	userIdValue, _ := ctx.Get("userId")
	userID := userIdValue.(uint)

	blog, err := c.service.RequestPublish(ctx, uint(blogID), userID)
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
	blogID, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	userIdValue, _ := ctx.Get("userId")
	adminID := userIdValue.(uint)

	var req ApproveBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	blog, err := c.service.ApproveBlog(ctx, uint(blogID), adminID, req)
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

func (c *BlogController) MyBlogs(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	userIdValue, _ := ctx.Get("userId")
	authorID := userIdValue.(uint)

	blogs, total, err := c.service.MyBlogs(ctx, authorID, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total": total,
		"blogs": blogs,
	})
}

// Optional public/my route wrapper
func (c *BlogController) MyBlogsPublic(ctx *gin.Context) {
	c.MyBlogs(ctx)
}

// Update Blog
func (c *BlogController) UpdateBlog(ctx *gin.Context) {
	blogID, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	userIdValue, _ := ctx.Get("userId")
	userID := userIdValue.(uint)

	var req UpdateBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	blog, err := c.service.UpdateBlog(ctx, uint(blogID), userID, req)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Blog updated successfully",
		"blog":    blog,
	})
}

// Delete blog
func (c *BlogController) DeleteBlog(ctx *gin.Context) {
	blogID, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	userIdValue, _ := ctx.Get("userId")
	authorID := userIdValue.(uint)

	if err := c.service.DeleteBlog(ctx, uint(blogID), authorID); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, DeleteBlogResponse{Message: "Blog deleted"})
}

func (c *BlogController) ListBlogs(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	search := ctx.Query("search")

	filters := map[string]interface{}{}
	if catID := ctx.Query("category_id"); catID != "" {
		if id, err := strconv.Atoi(catID); err == nil {
			filters["category_id"] = id
		}
	}
	if authorID := ctx.Query("author_id"); authorID != "" {
		if id, err := strconv.Atoi(authorID); err == nil {
			filters["author_id"] = id
		}
	}

	blogs, total, err := c.service.ListBlogs(ctx, filters, search, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total": total,
		"blogs": blogs,
	})
}

func (c *BlogController) GetBlogDetail(ctx *gin.Context) {
	blogID, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	blog, err := c.service.GetBlogDetail(ctx, uint(blogID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Blog not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, blog)
}

// Categories
func (c *BlogController) CreateCategory(ctx *gin.Context) {
	var req CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	cat, err := c.service.CreateCategory(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, CategoryResponse{
		ID:          cat.ID,
		Name:        cat.Name,
		Description: cat.Slug,
	})
}

func (c *BlogController) UpdateCategory(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	var req UpdateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	cat, err := c.service.UpdateCategory(ctx, uint(id), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, CategoryResponse{
		ID:          cat.ID,
		Name:        cat.Name,
		Description: cat.Slug,
	})
}

func (c *BlogController) DeleteCategory(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err := c.service.DeleteCategory(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, DeleteBlogResponse{
		Message: "Category deleted",
	})
}

func (c *BlogController) ListCategory(ctx *gin.Context) {
	cats, err := c.service.ListCategory(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	resp := make([]CategoryResponse, 0, len(cats))
	for _, cmodel := range cats {
		resp = append(resp, CategoryResponse{
			ID:          cmodel.ID,
			Name:        cmodel.Name,
			Description: cmodel.Slug,
		})
	}
	ctx.JSON(http.StatusOK, resp)
}
