package repository

import (
	"github.com/donnigundala/dg-core/contracts/foundation"
)

// Registry holds all registered repositories
type Registry struct {
	repositories []Registrable
}

// NewRegistry creates a new repository registry
func NewRegistry() *Registry {
	return &Registry{
		repositories: make([]Registrable, 0),
	}
}

// Register adds a repository to the registry
func (r *Registry) Register(repo Registrable) {
	r.repositories = append(r.repositories, repo)
}

// RegisterAll registers all repositories in the container
func (r *Registry) RegisterAll(app foundation.Application) error {
	for _, repo := range r.repositories {
		name := repo.Name()
		factory := repo.Factory

		// Register as singleton in the container
		app.Singleton(name, func() (interface{}, error) {
			return factory(app)
		})
	}

	return nil
}

// Count returns the number of registered repositories
func (r *Registry) Count() int {
	return len(r.repositories)
}

// GetNames returns all registered repository names
func (r *Registry) GetNames() []string {
	names := make([]string, 0, len(r.repositories))
	for _, repo := range r.repositories {
		names = append(names, repo.Name())
	}
	return names
}
