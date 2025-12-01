package main

import (
	"log"
	"skeleton-v2/bootstrap"
)

func main() {
	// Create a new application instance.
	app := bootstrap.NewApplication()

	// Boot the application.
	if err := app.Boot(); err != nil {
		log.Fatalf("Failed to boot application: %v", err)
	}

	// Start the application.
	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}
