package http

import (
	"skeleton/app/http/routes"

	"github.com/donnigundala/dg-core/foundation"
	coreHTTP "github.com/donnigundala/dg-core/http"
	"github.com/donnigundala/dg-core/http/health"
	"github.com/gin-gonic/gin"
)

// NewKernel configures the HTTP kernel (router, middleware, routes).
func NewKernel(app *foundation.Application) *gin.Engine {
	// Create Gin Engine using dg-core factory
	router := coreHTTP.NewRouter()

	// Register router instance in container
	app.Instance("router", router)

	// Setup health checks
	setupHealthChecks(router)

	// Apply global middleware
	router.Use(globalMiddleware()...)

	// Register application routes
	routes.Register(app, router)

	return router
}

func setupHealthChecks(router *gin.Engine) {
	healthManager := health.NewManager()
	healthManager.AddCheck(health.AlwaysHealthy("app"))

	router.GET("/health/live", health.LivenessHandler())
	router.GET("/health/ready", healthManager.ReadinessHandler())
	router.GET("/health", healthManager.HealthHandler())
}

func globalMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		coreHTTP.RequestIDWithDefault(),          // Request tracing (must be first)
		coreHTTP.LoggerWithDefault(),             // Logging with request ID
		coreHTTP.RecoveryWithDefault(),           // Panic recovery
		coreHTTP.CORSWithDefault(),               // CORS headers
		coreHTTP.SecurityHeadersWithDefault(),    // Security headers
		coreHTTP.BodySizeLimit(10 * 1024 * 1024), // 10MB limit
	}
}
