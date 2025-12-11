package routes

import (
	"skeleton/app/http/controllers"

	"github.com/donnigundala/dg-core/contracts/foundation"
	"github.com/gin-gonic/gin"
)

// Register registers the web routes.
func Register(app foundation.Application, router *gin.Engine) {
	// Welcome route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Skeleton V2!",
			"version": "2.0.0",
			"info":    "DG Framework Skeleton Application",
		})
	})

	// Helper to resolve userController lazily on each request
	getUserController := func() *controllers.UserController {
		userControllerInstance, err := app.Make("userController")
		if err != nil {
			panic("failed to resolve userController: " + err.Error())
		}
		return userControllerInstance.(*controllers.UserController)
	}

	// API group
	api := router.Group("/api/v1")
	{
		api.GET("/status", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "online",
				"app":    "skeleton-v2",
			})
		})

		// User routes
		api.POST("/users", func(c *gin.Context) {
			getUserController().Create(c)
		})
		api.GET("/users", func(c *gin.Context) {
			getUserController().List(c)
		})
		api.GET("/users/:id", func(c *gin.Context) {
			getUserController().Get(c)
		})
		api.PUT("/users/:id", func(c *gin.Context) {
			getUserController().Update(c)
		})
		api.DELETE("/users/:id", func(c *gin.Context) {
			getUserController().Delete(c)
		})
	}
}
