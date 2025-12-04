package providers

import (
	database "github.com/donnigundala/dg-database"
)

// NewDatabaseServiceProvider creates a database provider using the new plugin pattern.
// This is the recommended approach for most applications.
//
// The provider will automatically:
// - Resolve and integrate with the logger (if available)
// - Test database connection on boot
// - Handle graceful shutdown
//
// For advanced customization, you can still create a custom provider
// using database.NewManager() directly.
func NewDatabaseServiceProvider() *database.DatabaseServiceProvider {
	return &database.DatabaseServiceProvider{}
}
