package providers

import (
	"fmt"

	queue "github.com/donnigundala/dg-queue"
	"github.com/donnigundala/dg-queue/drivers/memory"
	"github.com/donnigundala/dg-queue/drivers/redis"
	goredis "github.com/redis/go-redis/v9"
)

// NewQueueServiceProvider creates a queue provider using the new plugin pattern.
// This is the recommended approach for most applications.
//
// The provider will automatically:
// - Resolve and integrate with the logger (if available)
// - Configure the driver based on config
// - Handle graceful shutdown
//
// For advanced customization, you can still create a custom provider
// using queue.New() and queue.SetDriver() directly.
func NewQueueServiceProvider(config queue.Config) *queue.QueueServiceProvider {
	return &queue.QueueServiceProvider{
		Config:        config,
		DriverFactory: createQueueDriver,
	}
}

// createQueueDriver creates a queue driver based on configuration.
// This function handles the driver-specific logic that can't be in the
// queue package due to import cycles.
func createQueueDriver(cfg queue.Config) (queue.Driver, error) {
	switch cfg.Driver {
	case "memory":
		return memory.NewDriver(), nil

	case "redis":
		// Parse Redis options from config
		redisOptions := &goredis.Options{
			Addr:     "localhost:6379", // Default
			Password: "",               // Default
			DB:       0,                // Default
		}

		// Override with config options if provided
		if addr, ok := cfg.Options["addr"].(string); ok {
			redisOptions.Addr = addr
		}
		if password, ok := cfg.Options["password"].(string); ok {
			redisOptions.Password = password
		}
		if db, ok := cfg.Options["db"].(int); ok {
			redisOptions.DB = db
		}

		prefix := cfg.Prefix
		if prefix == "" {
			prefix = "queue"
		}

		return redis.NewDriver(prefix, redisOptions)

	default:
		return nil, fmt.Errorf("unsupported queue driver: %s", cfg.Driver)
	}
}
