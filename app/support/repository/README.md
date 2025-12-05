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
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	GetByID(ctx context.Context, id uint) (*models.Product, error)
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

// ... implement other methods
```

### 2. Register in Loader

Add one line to `app/repositories/loader.go`:

```go
func LoadAll(app foundation.Application) error {
	registry := repository.NewRegistry()
	
	// Existing repositories
	registry.Register(repository.NewBaseRepository("userRepository", func(app foundation.Application) (interface{}, error) {
		db, err := getDB(app)
		if err != nil {
			return nil, err
		}
		return NewUserRepository(db), nil
	}))
	
	// Add your new repository here
	registry.Register(repository.NewBaseRepository("productRepository", func(app foundation.Application) (interface{}, error) {
		db, err := getDB(app)
		if err != nil {
			return nil, err
		}
		return NewProductRepository(db), nil
	}))
	
	return registry.RegisterAll(app)
}
```

**Note**: The `getDB()` helper function eliminates repetitive database resolution boilerplate.

That's it! No need to modify `app/providers/repository_provider.go`.

## Using Repositories in Services

Repositories are resolved from the container the same way as before:

```go
// app/services/user_service.go
package services

import (
	"skeleton-v2/app/repositories"
	"github.com/donnigundala/dg-core/contracts/foundation"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(app foundation.Application) (*UserService, error) {
	// Resolve from container
	repoInstance, err := app.Make("userRepository")
	if err != nil {
		return nil, err
	}
	
	return &UserService{
		userRepo: repoInstance.(repositories.UserRepository),
	}, nil
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
