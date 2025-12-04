package providers

import (
	scheduler "github.com/donnigundala/dg-scheduler"
)

// NewSchedulerServiceProvider creates a scheduler provider using the new plugin pattern.
// This is the recommended approach for most applications.
//
// The provider will automatically:
// - Resolve and integrate with the queue (if available)
// - Resolve and integrate with the logger (if available)
// - Start the scheduler on boot
// - Stop the scheduler on shutdown
//
// For advanced customization, you can still create a custom provider
// using scheduler.New() or scheduler.NewWithConfig() directly.
func NewSchedulerServiceProvider() *scheduler.SchedulerServiceProvider {
	return &scheduler.SchedulerServiceProvider{}
}
