package controllers

import (
	"skeleton/app/services"

	"github.com/donnigundala/dg-core/contracts/foundation"
	"github.com/donnigundala/dg-core/validation"
)

// Controllers holds all application controllers.
type Controllers struct {
	User *UserController
}

// Initialize creates and wires all controllers with their dependencies.
func Initialize(app foundation.Application) *Controllers {
	// Resolve shared dependencies
	validatorInstance, err := app.Make("validator")
	if err != nil {
		panic("failed to resolve validator: " + err.Error())
	}
	validator := validatorInstance.(*validation.Validator)

	// Resolve user service
	userServiceInstance, err := app.Make("userService")
	if err != nil {
		panic("failed to resolve user service: " + err.Error())
	}
	userService := userServiceInstance.(services.UserService)

	// Create and return all controllers
	return &Controllers{
		User: NewUserController(userService, validator),
	}
}
