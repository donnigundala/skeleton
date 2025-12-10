# Repository Infrastructure

This directory contains the infrastructure for repository auto-discovery. The actual repository definitions are in `app/repositories/`.

## Architecture

```
app/
├── repositories/              # Repository definitions (business logic)
│   ├── user_repository.go
│   ├── loader.go             # Registers all repositories
│   └── ...
└── support/repository/       # Infrastructure (reusable)
    ├── repository.go         # Registrable interface
    ├── registry.go           # Repository registry
    └── README.md             # This file
```

## Creating a New Repository

### 1. Create Repository Interface and Implementation

```go
// app/repositories/product_repository.go
package repositories

import (
	"context"
	"skeleton-v2/app/models"
	
	"github.com/donnigundala/dg-core/contracts/foundation"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	// ... other methods
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// MustResolveProductRepository resolves the product repository from the container.
func MustResolveProductRepository(app foundation.Application) ProductRepository {
	repo, err := app.Make("productRepository")
	if err != nil {
		panic("failed to resolve product repository: " + err.Error())
	}
	return repo.(ProductRepository)
}
```

### 2. Register in Loader

Add one line to `app/repositories/loader.go`:

```go
func LoadAll(app foundation.Application) error {
    // Resolve DB once
	dbManager := database.MustResolve(app)
	db := dbManager.DB()

	registry := repository.NewRegistry()
	
	// Existing repositories
	registry.Register(repository.NewBaseRepository("userRepository", func(app foundation.Application) (interface{}, error) {
		return NewUserRepository(db), nil
	}))
	
	// Add your new repository here
	registry.Register(repository.NewBaseRepository("productRepository", func(app foundation.Application) (interface{}, error) {
		return NewProductRepository(db), nil
	}))
	
	return registry.RegisterAll(app)
}
```

## Using Repositories in Services

Repositories are injected into services using the `MustResolveX` helper:

```go
// app/services/product_service.go
func NewProductService(app foundation.Application) ProductService {
    return &productService{
        repo: repositories.MustResolveProductRepository(app),
    }
}
```

## Benefits

- ✅ **Simple**: Add new repositories with one line in `loader.go`
- ✅ **Consistent**: Same pattern as scheduled jobs
- ✅ **Clean**: Separation of infrastructure and business logic
- ✅ **No provider changes**: All registration in one place
- ✅ **Backward compatible**: Existing code works unchanged

## Registrable Interface

The `Registrable` interface allows repositories to be auto-registered:

```go
type Registrable interface {
	Name() string                                           // Container binding name
	Factory(app foundation.Application) (interface{}, error) // Factory function
}
```

Use `repository.NewBaseRepository(name, factory)` for simple registration.
