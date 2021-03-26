package monitor

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/utkarsh-pro/heamon/models"
	"github.com/utkarsh-pro/heamon/pkg/eventbus"
)

// Monitor acts a wrapper for configuration
// and status and will attach eventbus to them
// for observing changes
type Monitor struct {
	Status *Status
	Config *Config

	Prober *Prober
}

// SetupSubscribers sets up all of the
// initial subscribers
func (m *Monitor) SetupSubscribers() {
	// Setup subscriber for config updates
	m.setupConfigUpdateSubscribers()

	// Setup subscriber for status updates
	m.setupStatusUpdateSubscribers()
}

// New returns a pointer to an instance of the monitor
func New(config models.Config) *Monitor {
	st := &Status{
		Status:        models.NewStatus(config.Monitor.Services),
		subscriberIds: map[string]int64{},
	}
	cfg := &Config{Config: &config}

	prober := NewProber(cfg)
	prober.Start()

	mon := &Monitor{
		Status: st,
		Config: cfg,

		Prober: prober,
	}

	mon.SetupSubscribers()

	return mon
}

func (m *Monitor) setupConfigUpdateSubscribers() {
	eventbus.Bus.Subscribe(ConfigUpdate, func(data ...interface{}) {
		if m.Prober != nil {
			// Terminate old prober
			m.Prober.Terminate()

			// Refresh status
			m.Status.Refresh(m.Config.Monitor.Services)

			// Setup subscribers again
			m.setupStatusUpdateSubscribers()

			// Create new Prober
			m.Prober = NewProber(m.Config)

			// Start new Prober
			m.Prober.Start()
		}
	})
}

func (m *Monitor) setupStatusUpdateSubscribers() {
	for _, svc := range m.Status.Report {
		name := svc.Name
		topic := fmt.Sprintf("%s.%s", StatusUpdate, name)

		m.Status.subscriberIds[topic] = eventbus.Bus.Subscribe(topic, func(data ...interface{}) {
			if len(data) != 1 {
				logrus.Error("invalid data received")
				return
			}
			stat, ok := data[0].(models.HealthStatus)
			if !ok {
				logrus.Error("invalid type conversion to models.HealthStatus")
				return
			}

			m.Status.UpdateStatus(name, stat)
		})
	}
}

// Stop stops the monitoring
func (m *Monitor) Stop() {
	m.Prober.Terminate()
}
