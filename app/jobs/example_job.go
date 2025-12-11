package jobs

import (
	"fmt"
	"log/slog"
	"time"

	"skeleton/app/support/scheduler"
)

// ExampleScheduledJob is a simple scheduled job that runs periodically
type ExampleScheduledJob struct {
	scheduler.BaseJob
	logger *slog.Logger
}

// NewExampleScheduledJob creates a new example scheduled job
func NewExampleScheduledJob(logger *slog.Logger) *ExampleScheduledJob {
	return &ExampleScheduledJob{
		BaseJob: scheduler.NewBaseJob("example-job", "* * * * *", true),
		logger:  logger,
	}
}

// Handle executes the job logic
func (j *ExampleScheduledJob) Handle() error {
	j.logger.Info("Example scheduled job executed",
		"job", j.Name(),
		"time", time.Now().Format(time.RFC3339),
		"message", "This job runs every minute!")
	return nil
}

// init registers this job automatically
func init() {
	// Jobs will be registered via the provider
	fmt.Println("Example scheduled job loaded")
}
