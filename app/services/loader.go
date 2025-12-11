package services

import (
	"skeleton/app/repositories"
	"skeleton/app/support/service"

	"github.com/donnigundala/dg-core/contracts/foundation"
	queue "github.com/donnigundala/dg-queue"
)

// LoadAll loads and registers all services
// Add new services here to make them discoverable
func LoadAll(app foundation.Application) error {
	registry := service.NewRegistry()

	// Register all services here
	// This loader ensures all services are bound to the container
	// so they can be injected into controllers.
	//
	// To add a new service:
	// 1. Create the service interface and implementation in this package
	// 2. Add a MustResolveX helper in the service file
	// 3. Register it here using NewBaseService

	// Register User Service
	registry.Register(service.NewBaseService("userService", func(app foundation.Application) (interface{}, error) {
		// Resolve dependencies using type-safe helpers
		userRepo := repositories.MustResolveUserRepository(app)
		queueManager := queue.MustResolve(app)

		return NewUserService(userRepo, app, queueManager), nil
	}))

	return registry.RegisterAll(app)
}
