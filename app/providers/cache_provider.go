package providers

import (
	cache "github.com/donnigundala/dg-cache"
	cacheMemory "github.com/donnigundala/dg-cache/drivers/memory"
	cacheRedis "github.com/donnigundala/dg-cache/drivers/redis"
)

// NewCacheServiceProvider creates a cache provider using the new plugin pattern.
// This is the recommended approach for most applications.
//
// The provider will automatically:
// - Register common drivers (memory, redis)
// - Handle graceful shutdown
//
// For advanced customization, you can still create a custom provider
// using cache.NewManager() and cache.RegisterDriver() directly.
func NewCacheServiceProvider(config cache.Config) *cache.CacheServiceProvider {
	return &cache.CacheServiceProvider{
		Config: config,
		DriverFactories: map[string]cache.DriverFactory{
			"memory": cacheMemory.NewDriver,
			"redis":  cacheRedis.NewDriver,
		},
	}
}
