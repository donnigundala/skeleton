package repositories

import (
	"skeleton/app/support/repository"

	"github.com/donnigundala/dg-core/contracts/foundation"
	database "github.com/donnigundala/dg-database"
)

// LoadAll loads and registers all repositories
// Add new repositories here to make them discoverable
func LoadAll(app foundation.Application) error {
	// Resolve database connection once for all repositories
	// Use type-safe helper
	dbManager := database.MustResolve(app)
	db := dbManager.DB()

	registry := repository.NewRegistry()

	// Register all repositories here
	// This loader ensures all repositories are bound to the container
	// so they can be injected into services.
	//
	// To add a new repository:
	// 1. Create the repository interface and implementation in this package
	// 2. Add a MustResolveX helper in the repository file
	// 3. Register it here using NewBaseRepository
	registry.Register(repository.NewBaseRepository("userRepository", func(app foundation.Application) (interface{}, error) {
		return NewUserRepository(db), nil
	}))

	return registry.RegisterAll(app)
}
