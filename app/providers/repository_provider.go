package providers

import (
	"skeleton-v2/app/repositories"

	"github.com/donnigundala/dg-core/contracts/foundation"
	database "github.com/donnigundala/dg-database"
)

// RepositoryServiceProvider registers all application repositories.
type RepositoryServiceProvider struct{}

// NewRepositoryServiceProvider creates a new RepositoryServiceProvider.
func NewRepositoryServiceProvider() *RepositoryServiceProvider {
	return &RepositoryServiceProvider{}
}

// Register binds repositories into the container.
func (p *RepositoryServiceProvider) Register(app foundation.Application) error {
	// Register User Repository
	app.Singleton("userRepository", func() interface{} {
		dbManagerInstance, err := app.Make("database")
		if err != nil {
			panic("failed to resolve database manager: " + err.Error())
		}
		dbManager := dbManagerInstance.(*database.Manager)
		return repositories.NewUserRepository(dbManager.DB())
	})

	return nil
}

// Boot boots the service provider.
func (p *RepositoryServiceProvider) Boot(app foundation.Application) error {
	// Nothing to boot for repositories
	return nil
}
