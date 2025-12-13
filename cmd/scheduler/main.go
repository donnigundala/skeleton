package main

import (
	"log"
	"skeleton/bootstrap"
)

func main() {
	// Create a new application instance in scheduler mode.
	app := bootstrap.NewApplication(bootstrap.ModeScheduler)

	// Boot the application.
	if err := app.Boot(); err != nil {
		log.Fatalf("Failed to boot scheduler application: %v", err)
	}

	// Start the application.
	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start scheduler application: %v", err)
	}
}
