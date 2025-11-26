package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"skeleton-v2/app/http/routes"
	"skeleton-v2/app/providers"

	cache "github.com/donnigundala/dg-cache"
	"github.com/donnigundala/dg-core/config"
	"github.com/donnigundala/dg-core/contracts/http"
	"github.com/donnigundala/dg-core/errors"
	"github.com/donnigundala/dg-core/foundation"
	coreHTTP "github.com/donnigundala/dg-core/http"
	"github.com/donnigundala/dg-core/http/health"
	"github.com/donnigundala/dg-core/logging"
	"github.com/donnigundala/dg-core/validation"
)

// AppConfig represents the application configuration with validation.
type AppConfig struct {
	Name  string `mapstructure:"name" validate:"required,min=3"`
	Env   string `mapstructure:"env" validate:"required,oneof=development staging production"`
	Debug bool   `mapstructure:"debug"`
	Port  int    `mapstructure:"port" validate:"required,gte=1,lte=65535"`
}

// Application represents the bootstrapped application.
type Application struct {
	foundation *foundation.Application
	logger     *logging.Logger
	config     AppConfig
	server     *coreHTTP.HTTPServer
}

// NewApplication creates and bootstraps a new application instance.
func NewApplication() *Application {
	basePath, _ := os.Getwd()
	app := foundation.New(basePath)

	// Initialize logger
	logger := logging.New(logging.Config{
		Level:      slog.LevelInfo,
		Output:     os.Stdout,
		JSONFormat: false,
		AddSource:  false,
	})
	logging.SetDefault(logger)
	app.Instance("logger", logger)

	// Load configuration files
	config.Load()

	// Load and validate app configuration
	var appConfig AppConfig
	if err := config.Inject("app", &appConfig); err != nil {
		logger.Error("Failed to load app configuration", "error", err)
		os.Exit(1)
	}

	// Validate configuration
	validator := validation.NewValidator()
	if err := validator.ValidateStruct(context.Background(), &appConfig); err != nil {
		if valErr, ok := err.(*validation.Error); ok {
			logger.Error("Configuration validation failed", "errors", valErr.Errors)
		} else {
			logger.Error("Configuration validation failed", "error", err)
		}
		os.Exit(1)
	}

	logger.Info("Configuration validated successfully")

	// Enable debug mode if configured
	if appConfig.Debug {
		logger = logging.New(logging.Config{
			Level:      slog.LevelDebug,
			Output:     os.Stdout,
			JSONFormat: false,
			AddSource:  true,
		})
		logging.SetDefault(logger)
		app.Instance("logger", logger)
	}

	logger.Info("Starting application",
		"name", appConfig.Name,
		"env", appConfig.Env,
		"debug", appConfig.Debug,
	)

	// Load cache configuration
	var cacheConfig cache.Config
	if err := config.Inject("cache", &cacheConfig); err != nil {
		logger.Error("Failed to load cache configuration", "error", err)
		os.Exit(1)
	}

	// Register service providers
	app.Register(providers.NewCacheServiceProvider(cacheConfig))
	app.Register(providers.NewDatabaseServiceProvider())

	// Setup HTTP router
	router := setupRouter(app, logger)

	// Create HTTP kernel
	kernel := coreHTTP.NewKernel(app, router)

	// Create HTTP server
	addr := fmt.Sprintf(":%d", appConfig.Port)
	server := coreHTTP.NewHTTPServer(coreHTTP.Config{Addr: addr}, kernel)

	// Register shutdown hooks
	app.RegisterShutdownHook(func() {
		logger.Info("Shutting down HTTP server...")
		ctx, cancel := context.WithTimeout(context.Background(), 30)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			logger.Error("HTTP server shutdown error", "error", err)
		}
	})

	app.RegisterShutdownHook(func() {
		logger.Info("Executing cleanup: Closing resources...")
	})

	logger.Info("HTTP server configured", "addr", addr)
	logger.Info("Health endpoints available:")
	logger.Info("  - Liveness:  http://localhost" + addr + "/health/live")
	logger.Info("  - Readiness: http://localhost" + addr + "/health/ready")
	logger.Info("  - Detailed:  http://localhost" + addr + "/health")

	return &Application{
		foundation: app,
		logger:     logger,
		config:     appConfig,
		server:     server,
	}
}

// setupRouter configures the HTTP router with middleware and routes.
func setupRouter(app *foundation.Application, logger *logging.Logger) http.Router {
	// Bind router to container
	app.Singleton("router", func() interface{} {
		return coreHTTP.NewRouter()
	})

	// Get router instance
	routerInstance, err := app.Make("router")
	if err != nil {
		wrappedErr := errors.Wrap(err, "failed to resolve router").
			WithCode("ROUTER_RESOLUTION_FAILED").
			WithStatus(500)
		logger.Error("Failed to resolve router", "error", wrappedErr, "code", wrappedErr.Code())
		os.Exit(1)
	}
	router := routerInstance.(http.Router)

	// Setup health checks
	healthManager := health.NewManager()
	healthManager.AddCheck(health.AlwaysHealthy("app"))
	healthManager.AddCheck(health.SimpleCheck("custom", func(ctx context.Context) error {
		// Add custom health check logic here
		// For example: check database connection, cache connection, etc.
		return nil
	}))

	// Register health check routes
	router.Get("/health/live", health.LivenessHandler())
	router.Get("/health/ready", healthManager.ReadinessHandler())
	router.Get("/health", healthManager.HealthHandler())

	// Apply global middleware
	logger.Info("Applying global middleware")
	router.Use(
		coreHTTP.LoggerWithDefault(),          // Request logging
		coreHTTP.RecoveryWithDefault(),        // Panic recovery
		coreHTTP.CORSWithDefault(),            // CORS
		coreHTTP.SecurityHeadersWithDefault(), // Security headers
		coreHTTP.BodySizeLimit(10*1024*1024),  // 10MB limit
	)

	// Register application routes
	routes.Register(router)

	logger.Debug("Routes registered successfully")

	return router
}

// Start starts the HTTP server.
func (a *Application) Start() error {
	a.logger.Info("Application started. Press Ctrl+C to shutdown.")

	// Start server in goroutine
	go func() {
		if err := a.server.Start(); err != nil {
			a.logger.Error("HTTP server error", "error", err)
		}
	}()

	// Wait for shutdown signal
	a.foundation.WaitForShutdown()

	a.logger.Info("Application stopped gracefully")
	return nil
}
