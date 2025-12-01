package providers

import (
	"skeleton-v2/app/http/controllers"
	"skeleton-v2/app/repositories"
	"skeleton-v2/app/services"

	"github.com/donnigundala/dg-cache"
	"github.com/donnigundala/dg-core/contracts/foundation"
	"github.com/donnigundala/dg-core/validation"
	"github.com/donnigundala/dg-database"
	"github.com/donnigundala/dg-queue"
)

// AppServiceProvider registers the application's core components.
type AppServiceProvider struct{}

// NewAppServiceProvider creates a new AppServiceProvider.
func NewAppServiceProvider() *AppServiceProvider {
	return &AppServiceProvider{}
}

// Register binds the application's services into the container.
func (p *AppServiceProvider) Register(app foundation.Application) error {
	// Register Validator as a singleton
	app.Singleton("validator", func() interface{} {
		return validation.NewValidator()
	})

	// Register Repositories as singletons
	app.Singleton("userRepository", func() interface{} {
		dbManagerInstance, err := app.Make("database")
		if err != nil {
			panic("failed to resolve database manager: " + err.Error())
		}
		dbManager := dbManagerInstance.(*database.Manager)
		return repositories.NewUserRepository(dbManager.DB())
	})

	// Register Services as singletons
	app.Singleton("userService", func() interface{} {
		userRepoInstance, err := app.Make("userRepository")
		if err != nil {
			panic("failed to resolve user repository: " + err.Error())
		}
		cacheManagerInstance, err := app.Make("cache")
		if err != nil {
			panic("failed to resolve cache manager: " + err.Error())
		}
		queueManagerInstance, err := app.Make("queue")
		if err != nil {
			panic("failed to resolve queue manager: " + err.Error())
		}

		return services.NewUserService(
			userRepoInstance.(repositories.UserRepository),
			cacheManagerInstance.(*cache.Manager),
			queueManagerInstance.(*queue.Manager),
		)
	})

	// Register Controllers as singletons
	app.Singleton("userController", func() interface{} {
		userServiceInstance, err := app.Make("userService")
		if err != nil {
			panic("failed to resolve user service: " + err.Error())
		}
		validatorInstance, err := app.Make("validator")
		if err != nil {
			panic("failed to resolve validator: " + err.Error())
		}

		return controllers.NewUserController(
			userServiceInstance.(services.UserService),
			validatorInstance.(*validation.Validator),
		)
	})

	return nil
}

// Boot boots the service provider.
func (p *AppServiceProvider) Boot(app foundation.Application) error {
	// Nothing to boot for this provider
	return nil
}
