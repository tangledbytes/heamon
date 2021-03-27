package monitor

import (
	"net/http"
	"path"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/utkarsh-pro/heamon/pkg/eventbus"
	"github.com/utkarsh-pro/heamon/pkg/store/status"
)

type ProbeBot struct {
	name       string
	host       string
	endpoint   string
	interval   int
	sendUpdate func(service string, status status.HealthStatus)

	eb *eventbus.EventBus
}

// NewProbeBot returns a pointer to an instance of ProbBot
func NewProbeBot(
	eb *eventbus.EventBus,
	interval int,
	name, host, endpoint string,
	sendUpdate func(service string, status status.HealthStatus),
) *ProbeBot {
	return &ProbeBot{
		name:       name,
		host:       host,
		endpoint:   endpoint,
		interval:   interval,
		sendUpdate: sendUpdate,

		eb: eb,
	}
}

// Start the probbot
func (p *ProbeBot) Start() {
	ticker := time.NewTicker(time.Duration(p.interval) * time.Minute)

	p.eb.Subscribe(string(TerminateProbeBot), func(data ...interface{}) {
		logrus.Info("received cancellation for", p.name)
		ticker.Stop()
	})

	// First probe should happen immediately
	logrus.Info("Probing ", p.host)
	p.SendUpdate(probe(p.host, p.endpoint))

	for range ticker.C {
		logrus.Info("Probing ", p.host)
		p.SendUpdate(probe(p.host, p.endpoint))
	}
}

// SendUpdate sends health update via the method
// provided to it by the prober
func (p *ProbeBot) SendUpdate(ok bool) {
	if ok {
		p.sendUpdate(p.name, status.HealthOK)
		return
	}

	p.sendUpdate(p.name, status.HealthFail)
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
