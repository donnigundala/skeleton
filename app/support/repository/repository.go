package repository

import "github.com/donnigundala/dg-core/contracts/foundation"

// Registrable defines the interface for repositories that can be auto-registered
type Registrable interface {
	// Name returns the container binding name for this repository
	// Example: "userRepository", "productRepository"
	Name() string

	// Factory creates and returns the repository instance
	// This function will be called by the container when the repository is resolved
	Factory(app foundation.Application) (interface{}, error)
}

// BaseRepository provides a simple implementation helper
type BaseRepository struct {
	name    string
	factory func(foundation.Application) (interface{}, error)
}

// NewBaseRepository creates a new base repository registration
func NewBaseRepository(name string, factory func(foundation.Application) (interface{}, error)) *BaseRepository {
	return &BaseRepository{
		name:    name,
		factory: factory,
	}
}

// Name returns the repository name
func (b *BaseRepository) Name() string {
	return b.name
}

// Factory returns the repository factory function
func (b *BaseRepository) Factory(app foundation.Application) (interface{}, error) {
	return b.factory(app)
}
