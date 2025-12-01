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
	"github.com/donnigundala/dg-queue"
)

// AppConfig represents the application configuration.
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

// NewApplication creates a new application instance but does not boot it.
func NewApplication() *Application {
	basePath, _ := os.Getwd()
	app := foundation.New(basePath)

	// The rest of the initialization will be handled by the Boot method.
	return &Application{
		foundation: app,
	}
}

// Boot initializes and bootstraps the application.
func (a *Application) Boot() error {
	// Initialize Logger
	a.logger = a.setupLogger()
	a.foundation.Instance("logger", a.logger)
	logging.SetDefault(a.logger)

	// Load and Validate Configuration
	if err := a.loadConfig(); err != nil {
		return err
	}

	// Update logger if debug mode is enabled
	if a.config.Debug {
		a.logger = a.setupLogger()
		a.foundation.Instance("logger", a.logger)
		logging.SetDefault(a.logger)
	}

	a.logger.Info("Starting application",
		"name", a.config.Name,
		"env", a.config.Env,
		"debug", a.config.Debug,
	)

	// Register Service Providers
	if err := a.registerProviders(); err != nil {
		return err
	}

	// Boot the core application, which in turn boots all registered providers.
	if err := a.foundation.Boot(); err != nil {
		return errors.Wrap(err, "failed to boot service providers")
	}

	// Setup HTTP Server (now that providers are booted)
	a.server = a.setupHTTPServer()

	// Register Shutdown Hooks
	a.registerShutdownHooks()

	a.logger.Info("Application bootstrapped successfully")
	return nil
}

// Start starts the HTTP server and waits for a shutdown signal.
func (a *Application) Start() error {
	a.logger.Info("Application started. Press Ctrl+C to shutdown.")

	go func() {
		if err := a.server.Start(); err != nil {
			a.logger.Error("HTTP server error", "error", err)
		}
	}()

	a.foundation.WaitForShutdown()

	a.logger.Info("Application stopped gracefully")
	return nil
}

func (a *Application) setupLogger() *logging.Logger {
	level := slog.LevelInfo
	addSource := false
	if a.config.Debug {
		level = slog.LevelDebug
		addSource = true
	}
	return logging.New(logging.Config{
		Level:      level,
		Output:     os.Stdout,
		JSONFormat: false,
		AddSource:  addSource,
	})
}

func (a *Application) loadConfig() error {
	if err := config.Load(); err != nil {
		return errors.Wrap(err, "failed to load configuration")
	}
	if err := config.Inject("app", &a.config); err != nil {
		return errors.Wrap(err, "failed to load app configuration")
	}

	validator := validation.NewValidator()
	if err := validator.ValidateStruct(context.Background(), &a.config); err != nil {
		return errors.Wrap(err, "configuration validation failed")
	}

	a.logger.Info("Configuration loaded and validated successfully")
	return nil
}

func (a *Application) registerProviders() error {
	var cacheConfig cache.Config
	if err := config.Inject("cache", &cacheConfig); err != nil {
		return errors.Wrap(err, "failed to load cache configuration")
	}

	var queueConfig queue.Config
	if err := config.Inject("queue", &queueConfig); err != nil {
		return errors.Wrap(err, "failed to load queue configuration")
	}

	providersToRegister := []foundation.ServiceProvider{
		providers.NewCacheServiceProvider(cacheConfig),
		providers.NewQueueServiceProvider(queueConfig),
		providers.NewDatabaseServiceProvider(),
		providers.NewAppServiceProvider(),
	}

	for _, provider := range providersToRegister {
		if err := a.foundation.Register(provider); err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to register %T", provider))
		}
	}

	a.logger.Info("Service providers registered")
	return nil
}

func (a *Application) setupHTTPServer() *coreHTTP.HTTPServer {
	router := a.setupRouter()
	kernel := coreHTTP.NewKernel(a.foundation, router)
	addr := fmt.Sprintf(":%d", a.config.Port)
	server := coreHTTP.NewHTTPServer(coreHTTP.Config{Addr: addr}, kernel)

	a.logger.Info("HTTP server configured", "addr", addr)
	return server
}

func (a *Application) setupRouter() http.Router {
	a.foundation.Singleton("router", func() interface{} {
		return coreHTTP.NewRouter()
	})

	routerInstance, _ := a.foundation.Make("router")
	router := routerInstance.(http.Router)

	// Setup health checks
	healthManager := health.NewManager()
	healthManager.AddCheck(health.AlwaysHealthy("app"))
	router.Get("/health/live", health.LivenessHandler())
	router.Get("/health/ready", healthManager.ReadinessHandler())
	router.Get("/health", healthManager.HealthHandler())

	// Apply global middleware
	router.Use(
		coreHTTP.LoggerWithDefault(),
		coreHTTP.RecoveryWithDefault(),
		coreHTTP.CORSWithDefault(),
		coreHTTP.SecurityHeadersWithDefault(),
		coreHTTP.BodySizeLimit(10*1024*1024),
	)

	// Register application routes
	routes.Register(a.foundation, router)

	a.logger.Debug("Routes registered successfully")
	return router
}

func (a *Application) registerShutdownHooks() {
	a.foundation.RegisterShutdownHook(func() {
		a.logger.Info("Shutting down HTTP server...")
		ctx, cancel := context.WithTimeout(context.Background(), 30)
		defer cancel()
		if err := a.server.Shutdown(ctx); err != nil {
			a.logger.Error("HTTP server shutdown error", "error", err)
		}
	})

	a.foundation.RegisterShutdownHook(func() {
		a.logger.Info("Executing cleanup: Closing resources...")
	})
}
