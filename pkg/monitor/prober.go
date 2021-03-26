package monitor

import "github.com/utkarsh-pro/heamon/pkg/eventbus"

type Prober struct {
	config *Config

	eb *eventbus.EventBus
}

// NewProber returns a pointer to an instance
// of the prober
func NewProber(config *Config) *Prober {
	return &Prober{
		config: config,
		eb:     eventbus.New(),
	}
}

// Start starts prober and its probebots
func (p *Prober) Start() {
	for _, svc := range p.config.Monitor.Services {
		// Create new Probe Bot and start it
		bot := NewProbeBot(p.eb, p.config.Monitor.Interval, svc.Name, svc.Host, svc.HealthCheckEndpoint, svc.Tolerance, p.config.Monitor.AverageOver)
		go bot.Start()
	}
}

// Terminates the prober
func (p *Prober) Terminate() {
	p.eb.Publish("prober.probebot.cancel")
}
