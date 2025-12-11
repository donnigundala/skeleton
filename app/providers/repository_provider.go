package providers

import (
	"skeleton/app/repositories"

	"github.com/donnigundala/dg-core/contracts/foundation"
)

// RepositoryServiceProvider registers all application repositories.
type RepositoryServiceProvider struct{}

// NewRepositoryServiceProvider creates a new RepositoryServiceProvider.
func NewRepositoryServiceProvider() *RepositoryServiceProvider {
	return &RepositoryServiceProvider{}
}

// Register binds repositories into the container using auto-discovery.
func (p *RepositoryServiceProvider) Register(app foundation.Application) error {
	// Auto-discover and register all repositories
	return repositories.LoadAll(app)
}

// Boot boots the service provider.
func (p *RepositoryServiceProvider) Boot(app foundation.Application) error {
	// Nothing to boot for repositories
	return nil
}
