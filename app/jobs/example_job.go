package jobs

import (
	"fmt"
	"time"

	queue "github.com/donnigundala/dg-queue"
)

// ExampleJobName is the name of the example job
const ExampleJobName = "example_job"

// ExampleJobHandler handles the example job
func ExampleJobHandler(job *queue.Job) error {
	fmt.Printf("[Job] Processing example job: %s\n", job.ID)

	// Simulate work
	time.Sleep(1 * time.Second)

	// Access payload
	if msg, ok := job.Payload.(string); ok {
		fmt.Printf("[Job] Payload: %s\n", msg)
	} else {
		fmt.Printf("[Job] Payload: %v\n", job.Payload)
	}

	fmt.Printf("[Job] Completed example job: %s\n", job.ID)
	return nil
}
