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
	dbManagerInstance, err := app.Make("database")
	if err != nil {
		return err
	}
	dbManager := dbManagerInstance.(*database.Manager)
	db := dbManager.DB()

	registry := repository.NewRegistry()

	// Register all repositories here
	registry.Register(repository.NewBaseRepository("userRepository", func(app foundation.Application) (interface{}, error) {
		return NewUserRepository(db), nil
	}))

	// Add more repositories here as needed:
	// registry.Register(repository.NewBaseRepository("productRepository", func(app foundation.Application) (interface{}, error) {
	//     return NewProductRepository(db), nil
	// }))

	return registry.RegisterAll(app)
}
