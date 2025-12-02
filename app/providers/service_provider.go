package providers

import (
	"skeleton-v2/app/repositories"
	"skeleton-v2/app/services"

	cache "github.com/donnigundala/dg-cache"
	"github.com/donnigundala/dg-core/contracts/foundation"
	queue "github.com/donnigundala/dg-queue"
)

// ServiceLayerProvider registers all application services.
type ServiceLayerProvider struct{}

// NewServiceLayerProvider creates a new ServiceLayerProvider.
func NewServiceLayerProvider() *ServiceLayerProvider {
	return &ServiceLayerProvider{}
}

// Register binds services into the container.
func (p *ServiceLayerProvider) Register(app foundation.Application) error {
	// Register User Service
	app.Singleton("userService", func() interface{} {
		userRepoInstance, err := app.Make("userRepository")
		if err != nil {
			panic("failed to resolve user repository: " + err.Error())
		}
		cacheManagerInstance, err := app.Make("cache")
		if err != nil {
			panic("failed to resolve cache manager: " + err.Error())
		}
		queueManagerInstance, err := app.Make("queue")
		if err != nil {
			panic("failed to resolve queue manager: " + err.Error())
		}

		return services.NewUserService(
			userRepoInstance.(repositories.UserRepository),
			cacheManagerInstance.(*cache.Manager),
			queueManagerInstance.(*queue.Manager),
		)
	})

	return nil
}

// Boot boots the service provider.
func (p *ServiceLayerProvider) Boot(app foundation.Application) error {
	// Nothing to boot for services
	return nil
}
