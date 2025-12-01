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

// Register only binds the key, but does not resolve or connect.
// The actual instance will be created in the Boot phase.
func (p *DatabaseServiceProvider) Register(app contractFoundation.Application) error {
	app.Singleton("database", func() interface{} {
		// This resolver will be called during the Boot phase,
		// or when "database" is requested for the first time after booting.
		var dbConfig database.Config
		if err := config.Inject("database", &dbConfig); err != nil {
			panic("failed to load database configuration: " + err.Error())
		}

		loggerInstance, _ := app.Make("logger")
		logger := loggerInstance.(database.Logger)

		manager, err := database.NewManager(dbConfig, logger)
		if err != nil {
			panic("failed to create and connect to database: " + err.Error())
		}
		return manager
	})
	return nil
}

// Boot does NOT force connection. Database will connect on first use (lazy loading).
// This prevents blocking startup if database is unavailable.
func (p *DatabaseServiceProvider) Boot(app contractFoundation.Application) error {
	// No eager connection - database manager will connect when first resolved
	return nil
}
