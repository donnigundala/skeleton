package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	appHTTP "skeleton/app/http"
	"skeleton/app/jobs"
	"skeleton/app/providers"

	cache "github.com/donnigundala/dg-cache"
	"github.com/donnigundala/dg-core/config"
	"github.com/donnigundala/dg-core/errors"
	"github.com/donnigundala/dg-core/foundation"
	coreHTTP "github.com/donnigundala/dg-core/http"
	"github.com/donnigundala/dg-core/logging"
	"github.com/donnigundala/dg-core/validation"
	filesystem "github.com/donnigundala/dg-filesystem"
	firebase "github.com/donnigundala/dg-firebase"

	// "github.com/donnigundala/dg-filesystem/drivers/s3" // Uncomment to enable S3 driver
	queue "github.com/donnigundala/dg-queue"
	scheduler "github.com/donnigundala/dg-scheduler"
)

// AppMode defines the application execution mode.
type AppMode string

const (
	// ModeWeb runs the application as an HTTP server.
	ModeWeb AppMode = "web"
	// ModeScheduler runs the application as a background job scheduler.
	ModeScheduler AppMode = "scheduler"
)

// AppConfig represents the application configuration.
type AppConfig struct {
	Name  string `mapstructure:"name" validate:"required,min=3"`
	Env   string `mapstructure:"env" validate:"required,oneof=development staging production"`
	Debug bool   `mapstructure:"debug"`
}

// Application represents the bootstrapped application.
type Application struct {
	foundation *foundation.Application
	logger     *logging.Logger
	config     AppConfig
	server     *coreHTTP.HTTPServer
	mode       AppMode
}

// NewApplication creates a new application instance with the specified mode.
func NewApplication(mode AppMode) *Application {
	basePath, _ := os.Getwd()
	app := foundation.New(basePath)

	return &Application{
		foundation: app,
		mode:       mode,
	}
}

// Boot initializes and bootstraps the application.
func (a *Application) Boot() error {
	// Initialize basic logger (before config is loaded)
	a.logger = a.setupLogger(false)
	a.foundation.Instance("logger", a.logger)
	logging.SetDefault(a.logger)

	// Load and Validate Configuration
	if err := a.loadConfig(); err != nil {
		return err
	}

	// Reconfigure logger with debug settings if enabled
	if a.config.Debug {
		a.logger = a.setupLogger(true)
		a.foundation.Instance("logger", a.logger)
		logging.SetDefault(a.logger)
	}

	a.logger.Info("Starting application",
		"name", a.config.Name,
		"env", a.config.Env,
		"debug", a.config.Debug,
		"mode", a.mode,
	)

	// Register Service Providers
	if err := a.registerProviders(); err != nil {
		return err
	}

	// Boot the core application, which in turn boots all registered providers.
	a.logger.Info("Booting service providers...")
	if err := a.foundation.Boot(); err != nil {
		return errors.Wrap(err, "failed to boot service providers")
	}
	a.logger.Info("Service providers booted successfully")

	// Setup HTTP Server (only in web mode)
	if a.mode == ModeWeb {
		a.server = a.setupHTTPServer()
	}

	// Register Shutdown Hooks
	a.registerShutdownHooks()

	a.logger.Info("Application bootstrapped successfully")
	return nil
}

// Start starts the application.
func (a *Application) Start() error {
	a.logger.Info("Application started. Press Ctrl+C to shutdown.")

	// Start all runnable services (Scheduler, Queue, etc.)
	// This uses the new lifecycle management system
	if err := a.foundation.StartServices(); err != nil {
		a.logger.Error("Failed to start services", "error", err)
		return err
	}

	switch a.mode {
	case ModeWeb:
		// Web mode: Start HTTP server
		go func() {
			if err := a.server.Start(); err != nil {
				a.logger.Error("HTTP server error", "error", err)
			}
		}()
	case ModeScheduler:
		// Scheduler mode: Register and run scheduled jobs
		if err := a.registerScheduledJobs(); err != nil {
			a.logger.Error("Failed to register scheduled jobs", "error", err)
			return err
		}
		a.logger.Info("Scheduler running in foreground (press Ctrl+C to stop)...")
	}

	a.foundation.WaitForShutdown()

	a.logger.Info("Application stopped gracefully")
	return nil
}

func (a *Application) setupLogger(debug bool) *logging.Logger {
	level := slog.LevelInfo
	addSource := false
	if debug {
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

	// Register providers in dependency order
	providersToRegister := []foundation.ServiceProvider{
		// Infrastructure layer
		providers.NewCacheServiceProvider(cacheConfig),
		providers.NewQueueServiceProvider(queueConfig),
		providers.NewSchedulerServiceProvider(), // Queue must be registered before Scheduler
		providers.NewDatabaseServiceProvider(),
		filesystem.NewFilesystemServiceProvider(), // Filesystem
		firebase.NewFirebaseServiceProvider(),     // Firebase integration

		// Application layer (order matters: Repositories → Services)
		providers.NewRepositoryServiceProvider(),
		providers.NewServiceLayerProvider(),
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
	a.logger.Info("Setting up HTTP server...")
	router := appHTTP.NewKernel(a.foundation)

	// Create Kernel with Gin Engine
	kernel := coreHTTP.NewKernel(a.foundation, router)

	// Load server configuration from server.yaml
	var serverConfig coreHTTP.Config
	if err := config.Inject("server", &serverConfig); err != nil {
		a.logger.Warn("Failed to load server config, using defaults", "error", err)
		serverConfig = coreHTTP.Config{
			Addr:         ":8080",
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		}
	}

	server := coreHTTP.NewHTTPServer(serverConfig, kernel, coreHTTP.WithHTTPLogger(a.logger.Underlying()))

	a.logger.Info("HTTP server configured", "addr", serverConfig.Addr)
	return server
}

func (a *Application) registerShutdownHooks() {
	a.foundation.RegisterShutdownHook(func() {
		if a.server != nil {
			a.logger.Info("Shutting down HTTP server...")
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if err := a.server.Shutdown(ctx); err != nil {
				a.logger.Error("HTTP server shutdown error", "error", err)
			}
		}
	})

	a.foundation.RegisterShutdownHook(func() {
		a.logger.Info("Executing cleanup: Closing resources...")
		a.foundation.StopServices()
	})
}

func (a *Application) registerScheduledJobs() error {
	// Get scheduler instance using type-safe helper
	// We use Resolve here since scheduler is optional/might fail
	scheduler, err := scheduler.Resolve(a.foundation)
	if err != nil {
		return fmt.Errorf("failed to resolve scheduler: %w", err)
	}

	// Load all jobs from the registry
	registry := jobs.LoadAll(a.foundation.Log())

	// Schedule all enabled jobs
	if err := registry.ScheduleAll(scheduler); err != nil {
		return err
	}

	return nil
}
