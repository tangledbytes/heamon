package monitor

import (
	"github.com/utkarsh-pro/heamon/pkg/eventbus"
	"github.com/utkarsh-pro/heamon/pkg/store"
)

type Prober struct {
	config store.Config
	status store.Status

	eb *eventbus.EventBus
}

// NewProber returns a pointer to an instance
// of the prober
func NewProber(config store.Config, status store.Status) *Prober {
	return &Prober{
		config: config,
		status: status,
		eb:     eventbus.New(),
	}
}

// Start starts prober and its probebots
func (p *Prober) Start() {
	// Get static copy of the config object
	cfg := p.config.Copy()

	for _, svc := range cfg.Monitor.Services {
		// Create new Probe Bot and start it
		go NewProbeBot(
			p.eb,
			cfg.Monitor.Interval,
			svc.Name,
			svc.Host,
			svc.HealthCheckEndpoint,
			p.status.Update,
		).Start()
	}
}

// Terminates the prober
func (p *Prober) Terminate() {
	p.eb.Publish(string(TerminateProbeBot))
}
