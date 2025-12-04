package main

import (
	"log"
	"skeleton-v2/bootstrap"
)

func main() {
	// Create a new scheduler application instance.
	app := bootstrap.NewSchedulerApplication()

	// Boot the application.
	if err := app.Boot(); err != nil {
		log.Fatalf("Failed to boot scheduler application: %v", err)
	}

	// Start the application.
	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start scheduler application: %v", err)
	}
}
