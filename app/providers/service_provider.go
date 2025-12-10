package providers

import (
	"skeleton-v2/app/repositories"
	"skeleton-v2/app/services"

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
		// Resolve dependencies using type-safe helpers
		userRepo := repositories.MustResolveUserRepository(app)
		queueManager := queue.MustResolve(app)

		return services.NewUserService(
			userRepo,
			app,
			queueManager,
		)
	})

	return nil
}

// Boot boots the service provider.
func (p *ServiceLayerProvider) Boot(app foundation.Application) error {
	// Nothing to boot for services
	return nil
}
