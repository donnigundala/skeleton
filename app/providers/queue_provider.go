package providers

import (
	"github.com/donnigundala/dg-core/contracts/foundation"
	queue "github.com/donnigundala/dg-queue"
	"github.com/donnigundala/dg-queue/drivers/memory"
)

// QueueServiceProvider registers the queue manager.
type QueueServiceProvider struct {
	Config queue.Config
}

// NewQueueServiceProvider creates a new QueueServiceProvider.
func NewQueueServiceProvider(config queue.Config) *QueueServiceProvider {
	return &QueueServiceProvider{Config: config}
}

// Register binds the queue manager into the container.
func (p *QueueServiceProvider) Register(app foundation.Application) error {
	app.Singleton("queue", func() interface{} {
		manager := queue.New(p.Config)

		// Register memory driver (for development/testing)
		if p.Config.Driver == "memory" {
			manager.SetDriver(memory.NewDriver())
		}

		return manager
	})
	return nil
}

// Boot boots the service provider.
func (p *QueueServiceProvider) Boot(app foundation.Application) error {
	// Nothing to boot for this provider
	return nil
}
