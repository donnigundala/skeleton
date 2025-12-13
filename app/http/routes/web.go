package routes

import (
	"skeleton/app/http/controllers"

	"github.com/donnigundala/dg-core/contracts/foundation"
	"github.com/gin-gonic/gin"
)

// Register registers the web routes.
func Register(app foundation.Application, router *gin.Engine) {
	// Initialize all controllers once
	ctrl := controllers.Initialize(app)

	// Welcome route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Skeleton V2!",
			"version": "2.0.0",
			"info":    "DG Framework Skeleton Application",
		})
	})

	// API group
	api := router.Group("/api/v1")
	{
		api.GET("/status", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "online",
				"app":    "skeleton-v2",
			})
		})

		// User routes - clean and direct
		api.POST("/users", ctrl.User.Create)
		api.GET("/users", ctrl.User.List)
		api.GET("/users/:id", ctrl.User.Get)
		api.PUT("/users/:id", ctrl.User.Update)
		api.DELETE("/users/:id", ctrl.User.Delete)
	}
}
