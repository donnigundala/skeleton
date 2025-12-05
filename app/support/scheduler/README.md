# Scheduled Jobs Infrastructure

This directory contains the infrastructure for scheduled jobs. The actual job definitions are in `app/jobs/`.

## Architecture

```
app/
├── jobs/                        # Job definitions (business logic)
│   ├── example_job.go
│   ├── loader.go                # Registers all jobs
│   └── ...
└── support/scheduler/           # Infrastructure (reusable)
    ├── job.go                   # ScheduledJob interface & BaseJob
    ├── registry.go              # Job registry
    └── README.md                # This file
```

## Creating a New Job

1. Create a new file in `app/jobs/` (e.g., `my_job.go`)
2. Implement the `ScheduledJob` interface:

```go
package jobs

import (
	"log/slog"
	"skeleton-v2/app/support/scheduler"
)

type MyJob struct {
	scheduler.BaseJob
	logger *slog.Logger
}

func NewMyJob(logger *slog.Logger) *MyJob {
	return &MyJob{
		BaseJob: scheduler.NewBaseJob("my-job", "0 * * * *", true),
		logger:  logger,
	}
}

func (j *MyJob) Handle() error {
	j.logger.Info("My job executed", "job", j.Name())
	// Your job logic here
	return nil
}
```

3. Register the job in `app/jobs/loader.go`:

```go
func LoadAll(logger *slog.Logger) *scheduler.Registry {
	registry := scheduler.NewRegistry(logger)
	
	registry.Register(NewExampleScheduledJob(logger))
	registry.Register(NewMyJob(logger))  // Add your new job here
	
	return registry
}
```

That's it! No need to modify `bootstrap/app.go`.

## Enabling/Disabling Jobs

To disable a job, pass `false` as the third parameter to `NewBaseJob`:

```go
BaseJob: scheduler.NewBaseJob("my-job", "0 * * * *", false)  // Disabled
```

## Cron Expression Examples

- `* * * * *` - Every minute
- `*/5 * * * *` - Every 5 minutes
- `0 * * * *` - Every hour
- `0 0 * * *` - Every day at midnight
- `0 9 * * 1` - Every Monday at 9 AM

## Job Interface

All jobs must implement the `ScheduledJob` interface:

```go
type ScheduledJob interface {
	Name() string           // Unique job name
	Schedule() string       // Cron expression
	Handle() error          // Job logic
	IsEnabled() bool        // Whether job should run
}
```

The `BaseJob` struct provides default implementations via `NewBaseJob(name, schedule, enabled)`.
