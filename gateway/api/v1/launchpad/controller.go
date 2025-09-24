package launchpad

import (
	"net/http"
	"strconv"

	"gateway/guards"
	"gateway/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LaunchpadController struct {
	service *LaunchpadService
}

func NewLaunchpadController(s *LaunchpadService) *LaunchpadController {
	return &LaunchpadController{service: s}
}

func (c *LaunchpadController) RegisterRoutes(r *gin.RouterGroup) {
	lpRoutes := r.Group("/launchpads")
	{
		lpRoutes.GET("/", c.GetAllLaunchpadHome)
		lpRoutes.GET("/:id", c.GetLaunchpadByID)
	}

	// authentication actions
	lpRoutesAuth := r.Group("/launchpads")
	lpRoutesAuth.Use(middlewares.AuthMiddleware())
	{
		// create & admin actions require permission
		lpRoutesAuth.POST("/", middlewares.RequirePermission(guards.PermCourseCRUD), c.CreateLaunchpad)
		lpRoutesAuth.POST("/:id/approve", middlewares.RequirePermission(guards.PermCourseCRUD), c.ApproveLaunchpad)

		// invest
		lpRoutesAuth.POST("/:id/invest", c.InvestInLaunchpad)

	}
}

func (c *LaunchpadController) CreateLaunchpad(ctx *gin.Context) {
	userIdValue, _ := ctx.Get("userId")

	userId, ok := userIdValue.(uint)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid userId type"})
		return
	}
	var dto CreateLaunchpadDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	lp, err := c.service.CreateLaunchpad(dto, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, lp)
}

func (c *LaunchpadController) GetLaunchpadByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid launchpad",
		})
		return
	}

	launchpad, err := c.service.GetLaunchpadByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "launchpad not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, launchpad)
}

func (c *LaunchpadController) GetAllLaunchpadHome(ctx *gin.Context) {
	groups, err := c.service.GetAllLaunchpadHome()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, groups)
}

// func (c *LaunchpadController) GetLaunchpads(ctx *gin.Context) {
// 	launchpads, err := c.service.GetLaunchpads()
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, launchpads)
// }

func (c *LaunchpadController) ApproveLaunchpad(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid launchpad",
		})
		return
	}

	launchpad, err := c.service.ApproveLaunchpad(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "launchpad not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, launchpad)
}

func (c *LaunchpadController) InvestInLaunchpad(ctx *gin.Context) {
	// Get Launchpad ID from URL
	idStr := ctx.Param("id")
	launchpadID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invailid launchpad id",
		})
		return
	}

	// Get UserID from auth middleware
	userIdValue, _ := ctx.Get("userId")
	userID, ok := userIdValue.(uint)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user id type",
		})
		return
	}

	// Get amount from request body
	var dto InvestDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Call service
	updatedLaunchpad, err := c.service.InvestInLaunchpad(userID, uint(launchpadID), dto.Amount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, updatedLaunchpad)
}
