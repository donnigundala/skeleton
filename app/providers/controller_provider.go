package providers

import (
	"skeleton-v2/app/http/controllers"
	"skeleton-v2/app/services"

	"github.com/donnigundala/dg-core/contracts/foundation"
	"github.com/donnigundala/dg-core/validation"
)

// ControllerServiceProvider registers all application controllers.
type ControllerServiceProvider struct{}

// NewControllerServiceProvider creates a new ControllerServiceProvider.
func NewControllerServiceProvider() *ControllerServiceProvider {
	return &ControllerServiceProvider{}
}

// Register binds controllers into the container.
func (p *ControllerServiceProvider) Register(app foundation.Application) error {
	// Register Validator (shared dependency)
	app.Singleton("validator", func() interface{} {
		return validation.NewValidator()
	})

	// Register User Controller
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
func (p *ControllerServiceProvider) Boot(app foundation.Application) error {
	// Nothing to boot for controllers
	return nil
}
