package wishlists

import (
	"net/http"
	"strconv"

	"gateway/guards"
	"gateway/middlewares"

	"github.com/gin-gonic/gin"
)

type WishlistController struct {
	service *WishlistService
}

func NewWishlistController(service *WishlistService) *WishlistController {
	return &WishlistController{service: service}
}

func (c *WishlistController) RegisterRoutes(r *gin.RouterGroup) {
	wishlistRoutes := r.Group("/wishlist")
	wishlistRoutes.Use(middlewares.AuthMiddleware())
	wishlistRoutes.Use(middlewares.UserValidatorMiddleware())
	{
		wishlistRoutes.POST("/", middlewares.RequirePermission(guards.PermWishlistAdd), c.AddToWishlist)
		wishlistRoutes.GET("/", middlewares.RequirePermission(guards.PermWishlistView), c.GetWishlist)
	}
}

func (c *WishlistController) AddToWishlist(ctx *gin.Context) {
	var dto AddToWishlistDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body" + err.Error(),
		})
		return
	}

	userId := ctx.MustGet("userId").(uint)

	wishlistItem, err := c.service.AddToWishlist(userId, dto.CourseID)
	if err != nil {
		if err.Error() == "course already in wishlist" {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err.Error() == "course not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add to wishlist",
		})
	}
	ctx.JSON(http.StatusCreated, wishlistItem)
}

func (c *WishlistController) GetWishlist(ctx *gin.Context) {
	userId := ctx.MustGet("userId").(uint)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	wishlist, err := c.service.GetWishlist(userId, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve wishlist",
		})
		return
	}

	ctx.JSON(http.StatusOK, wishlist)
}
