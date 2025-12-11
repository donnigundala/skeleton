package jobs

import (
	"log/slog"

	"skeleton/app/support/scheduler"
)

// LoadAll loads all available scheduled jobs
// Add new jobs here to make them discoverable
func LoadAll(logger *slog.Logger) *scheduler.Registry {
	registry := scheduler.NewRegistry(logger)

	// Register all jobs here
	// New jobs can be added simply by creating a new instance and registering it
	registry.Register(NewExampleScheduledJob(logger))

	// Add more jobs here as needed:
	// registry.Register(NewAnotherJob(logger))
	// registry.Register(NewYetAnotherJob(logger))

	return registry
}
