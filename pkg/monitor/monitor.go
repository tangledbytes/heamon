package monitor

import (
	"github.com/utkarsh-pro/heamon/pkg/store"
	"github.com/utkarsh-pro/heamon/pkg/store/config"
)

// Monitor acts a wrapper for configuration
// and status and will attach eventbus to them
// for observing changes
type Monitor struct {
	Status store.Status
	Config store.Config

	Prober *Prober
}

// New returns a pointer to an instance of the monitor
func New(Config store.Config, Status store.Status) *Monitor {
	prober := NewProber(Config, Status)
	prober.Start()

	mon := &Monitor{
		Status: Status,
		Config: Config,
		Prober: prober,
	}

	mon.setupConfigWatcher()

	return mon
}

// setupConfigWatcher sets up all of the
// initial subscribers
func (m *Monitor) setupConfigWatcher() {
	m.Config.Watch(config.UPDATE, func(*config.Config) {
		if m.Prober != nil {
			// Terminate old prober
			m.Prober.Terminate()

			// Create new Prober
			m.Prober = NewProber(m.Config, m.Status)

			// Start new Prober
			m.Prober.Start()
		}
	})
}

// Stop stops the monitoring
func (m *Monitor) Stop() {
	m.Prober.Terminate()
}
