# Service Infrastructure

This directory contains the infrastructure for service auto-discovery. The actual service definitions are in `app/services/`.

## Architecture

```
app/
├── services/                 # Service definitions (business logic)
│   ├── user_service.go
│   ├── loader.go             # Registers all services
│   └── ...
└── support/service/          # Infrastructure (reusable)
    ├── service.go            # Registrable interface
    ├── registry.go           # Service registry
    └── README.md             # This file
```

## Creating a New Service

### 1. Create Service Interface, Implementation, and Helper

```go
// app/services/product_service.go
package services

import (
	"context"
	"skeleton-v2/app/repositories"
	"github.com/donnigundala/dg-core/contracts/foundation"
)

type ProductService interface {
	Create(ctx context.Context, name string) error
}

type productService struct {
	repo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) Create(ctx context.Context, name string) error {
	// Business logic
	return nil
}

// MustResolveProductService resolves the product service from the container.
// This helper makes injection easy and type-safe.
func MustResolveProductService(app foundation.Application) ProductService {
	svc, err := app.Make("productService")
	if err != nil {
		panic("failed to resolve product service: " + err.Error())
	}
	return svc.(ProductService)
}
```

### 2. Register in Loader

Add one block to `app/services/loader.go`:

```go
func LoadAll(app foundation.Application) error {
	registry := service.NewRegistry()
	
	// ... existing services ...
	
	// Register Product Service
	registry.Register(service.NewBaseService("productService", func(app foundation.Application) (interface{}, error) {
		// Resolve dependencies using type-safe helpers
		productRepo := repositories.MustResolveProductRepository(app)
		
		return NewProductService(productRepo), nil
	}))
	
	return registry.RegisterAll(app)
}
```

## Using Services in Controllers

Services are injected into controllers using the `MustResolveX` helper:

```go
// app/http/controllers/product_controller.go
func NewProductController(app foundation.Application) *ProductController {
    return &ProductController{
        service: services.MustResolveProductService(app),
    }
}
```

## Benefits

- ✅ **Simple**: Add new services with standard pattern
- ✅ **Type-Safe**: Usage of helpers prevents interface casting errors at runtime
- ✅ **Clean**: Separation of infrastructure and business logic
- ✅ **No provider changes**: All registration in one place
