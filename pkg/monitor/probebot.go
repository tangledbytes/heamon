package monitor

import (
	"net/http"
	"path"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/utkarsh-pro/heamon/pkg/eventbus"
	"github.com/utkarsh-pro/heamon/pkg/store/config"
	"github.com/utkarsh-pro/heamon/pkg/store/status"
)

type ProbeBot struct {
	name            string
	host            string
	endpoint        string
	interval        float64
	failureRate     float64
	degradeRate     float64
	initialDownTime float64
	sendUpdate      func(service string, status status.HealthStatus)

	totalProbeFails float64
	totalProbes     float64
	probeStartTime  time.Time
	eb              *eventbus.EventBus
}

// NewProbeBot returns a pointer to an instance of ProbBot
func NewProbeBot(
	eb *eventbus.EventBus,
	svc config.Service,
	sendUpdate func(service string, status status.HealthStatus),
) *ProbeBot {
	return &ProbeBot{
		name:            svc.Name,
		host:            svc.Host,
		endpoint:        svc.HealthCheckEndpoint,
		interval:        svc.Interval,
		failureRate:     svc.Failure,
		degradeRate:     svc.Degraded,
		initialDownTime: svc.InitialDownTime,
		sendUpdate:      sendUpdate,

		totalProbeFails: 0,
		totalProbes:     0,
		probeStartTime:  time.Now(),
		eb:              eb,
	}
}

// Start the probbot
func (p *ProbeBot) Start() {
	ticker := time.NewTicker(time.Duration(p.interval) * time.Minute)

	p.eb.Subscribe(string(TerminateProbeBot), func(data ...interface{}) {
		logrus.Info("received cancellation for: ", p.name)
		ticker.Stop()
	})

	// First probe should happen immediately
	logrus.Info("Probing ", p.host)
	p.sendUpdate(p.name, p.analyze(probe(p.host, p.endpoint)))

	for range ticker.C {
		logrus.Info("Probing ", p.host)
		p.sendUpdate(p.name, p.analyze(probe(p.host, p.endpoint)))
	}
}

func (p *ProbeBot) analyze(ok bool) status.HealthStatus {
	p.totalProbes++
	if !ok {
		p.totalProbeFails++

		cr := (p.totalProbeFails / p.totalProbes) * 100

		if time.Since(p.probeStartTime.Add(time.Duration(p.initialDownTime)*time.Minute)) < 0 {
			return status.HealthUnknown
		}

		if cr > p.failureRate {
			return status.HealthFail
		}

		if cr > p.degradeRate {
			return status.HealthDegraded
		}
	}

	return status.HealthOK
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
