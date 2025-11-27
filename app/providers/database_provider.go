package providers

import (
	contractFoundation "github.com/donnigundala/dg-core/contracts/foundation"
	database "github.com/donnigundala/dg-database"
)

// DatabaseServiceProvider handles database initialization.
type DatabaseServiceProvider struct {
	Config database.Config `config:"database"`
}

// NewDatabaseServiceProvider creates a new database service provider.
func NewDatabaseServiceProvider() *DatabaseServiceProvider {
	return &DatabaseServiceProvider{}
}

// Register registers database services in the container.
func (p *DatabaseServiceProvider) Register(app contractFoundation.Application) error {
	app.Singleton("database", func() interface{} {
		// Config already injected by framework!

		// Get logger from container
		loggerInstance, err := app.Make("logger")
		if err != nil {
			panic(err)
		}

		logger, ok := loggerInstance.(database.Logger)
		if !ok {
			panic("logger does not implement database.Logger interface")
		}

		manager, err := database.NewManager(p.Config, logger)
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
