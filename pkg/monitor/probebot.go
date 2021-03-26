package monitor

import (
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/utkarsh-pro/heamon/models"
	"github.com/utkarsh-pro/heamon/pkg/eventbus"
)

type ProbeBot struct {
	name        string
	host        string
	endpoint    string
	interval    int
	tolerance   float64
	averageOver int

	eb *eventbus.EventBus
}

// NewProbeBot returns a pointer to an instance of ProbBot
func NewProbeBot(eb *eventbus.EventBus, interval int, name, host, endpoint string, tolerance float64, averageOver int) *ProbeBot {
	return &ProbeBot{
		name:        name,
		host:        host,
		endpoint:    endpoint,
		interval:    interval,
		tolerance:   tolerance,
		averageOver: averageOver,

		eb: eb,
	}
}

// Start the probbot
func (p *ProbeBot) Start() {
	ticker := time.NewTicker(time.Duration(p.interval) * time.Second)

	total := float64(1)
	failCount := float64(0)
	alerted := false

	p.eb.Subscribe("prober.probebot.cancel", func(data ...interface{}) {
		logrus.Info("received cancellation for", p.name)
		ticker.Stop()
	})

	// First prbe should happen immediately
	logrus.Info("Probing ", p.host)
	probeWrapper(p.name, p.host, p.endpoint)

	for range ticker.C {
		logrus.Info("Probing ", p.host)

		total++
		if !probeWrapper(p.name, p.host, p.endpoint) {
			failCount++
		}

		if !alerted && total > float64(p.averageOver) && failCount/total > p.tolerance {
			eventbus.Bus.Publish(PluginAlert, (failCount/total)*100, p.name)
			alerted = true
		}
	}
}

func probeWrapper(name, host, endpoint string) bool {
	topic := fmt.Sprintf("%s.%s", StatusUpdate, name)

	if probe(host, endpoint) {
		eventbus.Bus.Publish(topic, models.HealthOK)
		return true
	}

	eventbus.Bus.Publish(topic, models.HealthFail)
	return false
}

func probe(host, endpoint string) bool {
	res, err := http.Get("http://" + path.Join(host, endpoint))
	if err != nil {
		return false
	}

	if res.StatusCode != http.StatusOK {
		return false
	}

	return true
}
