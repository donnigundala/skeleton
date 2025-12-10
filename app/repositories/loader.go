package repositories

import (
	"skeleton-v2/app/support/repository"

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
	registry.Register(repository.NewBaseRepository("userRepository", func(app foundation.Application) (interface{}, error) {
		return NewUserRepository(db), nil
	}))

	return registry.RegisterAll(app)
}

// MustResolveUserRepository resolves the user repository from the container.
func MustResolveUserRepository(app foundation.Application) UserRepository {
	repo, err := app.Make("userRepository")
	if err != nil {
		panic("failed to resolve user repository: " + err.Error())
	}
	return repo.(UserRepository)
}
