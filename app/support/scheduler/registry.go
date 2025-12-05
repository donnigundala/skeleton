package scheduler

import (
	"fmt"
	"log/slog"
)

// Registry holds all registered scheduled jobs
type Registry struct {
	jobs   []ScheduledJob
	logger *slog.Logger
}

// NewRegistry creates a new job registry
func NewRegistry(logger *slog.Logger) *Registry {
	return &Registry{
		jobs:   make([]ScheduledJob, 0),
		logger: logger,
	}
}

// Register adds a job to the registry
func (r *Registry) Register(job ScheduledJob) {
	r.jobs = append(r.jobs, job)
	r.logger.Debug("Job registered", "name", job.Name(), "schedule", job.Schedule(), "enabled", job.IsEnabled())
}

// GetEnabledJobs returns all enabled jobs
func (r *Registry) GetEnabledJobs() []ScheduledJob {
	enabled := make([]ScheduledJob, 0)
	for _, job := range r.jobs {
		if job.IsEnabled() {
			enabled = append(enabled, job)
		}
	}
	return enabled
}

// ScheduleAll schedules all enabled jobs with the provided scheduler
func (r *Registry) ScheduleAll(scheduler interface {
	Schedule(cronExpr, name string, handler func() error) error
}) error {
	enabledJobs := r.GetEnabledJobs()

	for _, job := range enabledJobs {
		if err := scheduler.Schedule(job.Schedule(), job.Name(), job.Handle); err != nil {
			return fmt.Errorf("failed to schedule job '%s': %w", job.Name(), err)
		}
		r.logger.Info("Job scheduled", "name", job.Name(), "schedule", job.Schedule())
	}

	r.logger.Info("All jobs scheduled", "total", len(r.jobs), "enabled", len(enabledJobs))
	return nil
}
