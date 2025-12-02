package providers

import (
	"context"
	"fmt"
	"time"

	"github.com/donnigundala/dg-core/contracts/foundation"
	queue "github.com/donnigundala/dg-queue"
	"github.com/donnigundala/dg-queue/drivers/memory"
	"github.com/donnigundala/dg-queue/drivers/redis"
	goredis "github.com/redis/go-redis/v9"
)

// QueueServiceProvider registers the queue manager.
type QueueServiceProvider struct {
	Config queue.Config
}

// NewQueueServiceProvider creates a new QueueServiceProvider.
func NewQueueServiceProvider(config queue.Config) *QueueServiceProvider {
	return &QueueServiceProvider{Config: config}
}

// Register binds the queue manager into the container.
func (p *QueueServiceProvider) Register(app foundation.Application) error {
	app.Singleton("queue", func() interface{} {
		manager := queue.New(p.Config)

		// Configure driver based on config
		switch p.Config.Driver {
		case "memory":
			manager.SetDriver(memory.NewDriver())
		case "redis":
			// Parse Redis options from config
			redisOptions := &goredis.Options{
				Addr:     "localhost:6379", // Default
				Password: "",               // Default
				DB:       0,                // Default
			}

			// Override with config options if provided
			if addr, ok := p.Config.Options["addr"].(string); ok {
				redisOptions.Addr = addr
			}
			if password, ok := p.Config.Options["password"].(string); ok {
				redisOptions.Password = password
			}
			if db, ok := p.Config.Options["db"].(int); ok {
				redisOptions.DB = db
			}

			prefix := "queue"
			if pfx, ok := p.Config.Options["prefix"].(string); ok {
				prefix = pfx
			}

			redisDriver, err := redis.NewDriver(prefix, redisOptions)
			if err != nil {
				panic(fmt.Sprintf("failed to create Redis queue driver: %v", err))
			}
			manager.SetDriver(redisDriver)
		default:
			panic(fmt.Sprintf("unsupported queue driver: %s", p.Config.Driver))
		}

		return manager
	})
	return nil
}

// Boot boots the service provider.
func (p *QueueServiceProvider) Boot(app foundation.Application) error {
	// Nothing to boot for this provider
	return nil
}

// Shutdown gracefully stops the queue manager.
func (p *QueueServiceProvider) Shutdown(app foundation.Application) error {
	queueInstance, err := app.Make("queue")
	if err != nil {
		return nil // Queue not initialized, nothing to shutdown
	}

	manager := queueInstance.(*queue.Manager)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return manager.Stop(ctx)
}
