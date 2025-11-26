package providers

import (
	"github.com/donnigundala/dg-core/config"
	contractFoundation "github.com/donnigundala/dg-core/contracts/foundation"
	database "github.com/donnigundala/dg-database"
)

// DatabaseServiceProvider handles database initialization.
type DatabaseServiceProvider struct{}

// NewDatabaseServiceProvider creates a new database service provider.
func NewDatabaseServiceProvider() *DatabaseServiceProvider {
	return &DatabaseServiceProvider{}
}

// Register registers database services in the container.
func (p *DatabaseServiceProvider) Register(app contractFoundation.Application) error {
	app.Singleton("database", func() interface{} {
		var dbConfig database.Config
		if err := config.Inject("database", &dbConfig); err != nil {
			panic(err)
		}

		// Get logger from container
		loggerInstance, err := app.Make("logger")
		if err != nil {
			panic(err)
		}

		manager, err := database.NewManager(dbConfig, loggerInstance)
		if err != nil {
			panic(err)
		}

		return manager
	})
	return nil
}

// Boot boots the service provider.
func (p *DatabaseServiceProvider) Boot(app contractFoundation.Application) error {
	// Nothing to boot
	return nil
}
