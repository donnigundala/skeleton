package providers

import (
	"skeleton-v2/app/services"

	"github.com/donnigundala/dg-core/contracts/foundation"
)

// ServiceLayerProvider registers all application services.
type ServiceLayerProvider struct{}

// NewServiceLayerProvider creates a new ServiceLayerProvider.
func NewServiceLayerProvider() *ServiceLayerProvider {
	return &ServiceLayerProvider{}
}

// Register binds services into the container.
func (p *ServiceLayerProvider) Register(app foundation.Application) error {
	// Auto-discover and register all services
	return services.LoadAll(app)
}

// Boot boots the service provider.
func (p *ServiceLayerProvider) Boot(app foundation.Application) error {
	// Nothing to boot for services
	return nil
}
