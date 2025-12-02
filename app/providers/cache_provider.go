package providers

import (
	"log/slog"

	cache "github.com/donnigundala/dg-cache"
	cacheMemory "github.com/donnigundala/dg-cache/drivers/memory"
	contractFoundation "github.com/donnigundala/dg-core/contracts/foundation"
)

// CacheServiceProvider handles cache initialization.
type CacheServiceProvider struct {
	config cache.Config
}

// NewCacheServiceProvider creates a new cache service provider.
func NewCacheServiceProvider(config cache.Config) *CacheServiceProvider {
	return &CacheServiceProvider{
		config: config,
	}
}

// Register registers cache services in the container.
func (p *CacheServiceProvider) Register(app contractFoundation.Application) error {
	app.Singleton("cache", func() interface{} {
		manager, err := cache.NewManager(p.config)
		if err != nil {
			panic(err)
		}

		// Register drivers
		manager.RegisterDriver("memory", cacheMemory.NewDriver)

		return manager
	})
	return nil
}

// Boot boots the service provider.
func (p *CacheServiceProvider) Boot(app contractFoundation.Application) error {
	// Nothing to boot
	return nil
}

// Shutdown gracefully closes cache connections.
func (p *CacheServiceProvider) Shutdown(app contractFoundation.Application) error {
	cacheInstance, err := app.Make("cache")
	if err != nil {
		return nil // Cache not initialized, nothing to shutdown
	}

	slog.Info("Shutting down cache manager...")
	manager := cacheInstance.(*cache.Manager)
	return manager.Close()
}
