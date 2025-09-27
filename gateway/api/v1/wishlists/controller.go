package wishlists

import (
	"net/http"

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
	wishlistRoutes.Use(middlewares.RequirePermission(guards.PermWishlistAdd))
	{
		wishlistRoutes.POST("/", c.AddToWishlist)
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
