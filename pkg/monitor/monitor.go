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
func New(storeManager *store.Manager) *Monitor {
	st := storeManager.Status()
	cfg := storeManager.Config()

	prober := NewProber(cfg, st)
	prober.Start()

	mon := &Monitor{
		Status: st,
		Config: cfg,
		Prober: prober,
	}

	mon.SetupSubscribers()

	return mon
}

// SetupSubscribers sets up all of the
// initial subscribers
func (m *Monitor) SetupSubscribers() {
	m.Config.Watch(config.UPDATE, func() {
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
