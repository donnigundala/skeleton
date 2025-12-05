package scheduler

// ScheduledJob defines the interface that all scheduled jobs must implement
type ScheduledJob interface {
	// Name returns the unique name of the job
	Name() string

	// Schedule returns the cron expression for this job
	// Examples: "* * * * *" (every minute), "0 * * * *" (every hour)
	Schedule() string

	// Handle executes the job logic
	Handle() error

	// IsEnabled returns whether this job should be registered
	// This can be overridden by configuration
	IsEnabled() bool
}

// BaseJob provides default implementations for common job methods
type BaseJob struct {
	name     string
	schedule string
	enabled  bool
}

// NewBaseJob creates a new base job with the given parameters
func NewBaseJob(name, schedule string, enabled bool) BaseJob {
	return BaseJob{
		name:     name,
		schedule: schedule,
		enabled:  enabled,
	}
}

// Name returns the job name
func (b *BaseJob) Name() string {
	return b.name
}

// Schedule returns the cron schedule
func (b *BaseJob) Schedule() string {
	return b.schedule
}

// IsEnabled returns whether the job is enabled
func (b *BaseJob) IsEnabled() bool {
	return b.enabled
}
