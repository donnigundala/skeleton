package main

import (
	"log"
	"skeleton/bootstrap"
)

func main() {
	// Create a new application instance.
	app := bootstrap.NewApplication(bootstrap.ModeWeb)

	// Boot the application.
	if err := app.Boot(); err != nil {
		log.Fatalf("Failed to boot application: %v", err)
	}

	// Start the application.
	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}
