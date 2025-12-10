package service

import "github.com/donnigundala/dg-core/contracts/foundation"

// Registrable defines the interface for services that can be auto-registered
type Registrable interface {
	// Name returns the container binding name for this service
	// Example: "userService", "productService"
	Name() string

	// Factory creates and returns the service instance
	// This function will be called by the container when the service is resolved
	Factory(app foundation.Application) (interface{}, error)
}

// BaseService provides a simple implementation helper
type BaseService struct {
	name    string
	factory func(foundation.Application) (interface{}, error)
}

// NewBaseService creates a new base service registration
func NewBaseService(name string, factory func(foundation.Application) (interface{}, error)) *BaseService {
	return &BaseService{
		name:    name,
		factory: factory,
	}
}

// Name returns the service name
func (b *BaseService) Name() string {
	return b.name
}

// Factory returns the service factory function
func (b *BaseService) Factory(app foundation.Application) (interface{}, error) {
	return b.factory(app)
}
