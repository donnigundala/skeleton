package facades

import (
	"sync"

	"github.com/donnigundala/dg-core/foundation"
	scheduler "github.com/donnigundala/dg-scheduler"
)

var (
	schedulerInstance *scheduler.Scheduler
	schedulerOnce     sync.Once
)

// Scheduler provides a static-like interface to the scheduler service.
type Scheduler struct{}

// SetInstance sets the underlying scheduler instance.
// This is called by the service provider during registration/boot.
func SetScheduler(s *scheduler.Scheduler) {
	schedulerOnce.Do(func() {
		schedulerInstance = s
	})
}

// Resolve resolves the scheduler instance from the application container.
func ResolveScheduler(app *foundation.Application) {
	if instance, err := app.Make("scheduler"); err == nil {
		if s, ok := instance.(*scheduler.Scheduler); ok {
			SetScheduler(s)
		}
	}
}

// Start starts the scheduler.
func (s Scheduler) Start() error {
	if schedulerInstance != nil {
		return schedulerInstance.Start()
	}
	return nil
}

// Schedule schedules a job.
func (s Scheduler) Schedule(spec, name string, job func() error) error {
	if schedulerInstance != nil {
		return schedulerInstance.Schedule(spec, name, job)
	}
	return nil
}

// Stop stops the scheduler.
func (s Scheduler) Stop() {
	if schedulerInstance != nil {
		schedulerInstance.Stop()
	}
}
