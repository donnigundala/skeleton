package service

import (
	"github.com/donnigundala/dg-core/contracts/foundation"
)

// Registry holds all registered services
type Registry struct {
	services []Registrable
}

// NewRegistry creates a new service registry
func NewRegistry() *Registry {
	return &Registry{
		services: make([]Registrable, 0),
	}
}

// Register adds a service to the registry
func (r *Registry) Register(svc Registrable) {
	r.services = append(r.services, svc)
}

// RegisterAll registers all services in the container
func (r *Registry) RegisterAll(app foundation.Application) error {
	for _, svc := range r.services {
		name := svc.Name()
		factory := svc.Factory

		// Register as singleton in the container
		app.Singleton(name, func() (interface{}, error) {
			return factory(app)
		})
	}

	return nil
}

// Count returns the number of registered services
func (r *Registry) Count() int {
	return len(r.services)
}

// GetNames returns all registered service names
func (r *Registry) GetNames() []string {
	names := make([]string, 0, len(r.services))
	for _, svc := range r.services {
		names = append(names, svc.Name())
	}
	return names
}
