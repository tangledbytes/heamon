package status

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/utkarsh-pro/heamon/pkg/eventbus"
	"github.com/utkarsh-pro/heamon/pkg/store/config"
)

// Status stores the status of the services
type Status struct {
	UpdatedAt time.Time       `json:"updated_at,omitempty"`
	Report    []ServiceHealth `json:"report,omitempty"`

	eb *eventbus.EventBus `json:"-"`
	mu sync.Mutex         `json:"-"`
}

// ServiceHealth captures the metadata about the service
// along with its status
type ServiceHealth struct {
	config.Service `json:"service,omitempty"`
	HealthStatus   HealthStatus `json:"health_status,omitempty"`
}

// New returns a new instance of the status
func New(svcs []config.Service) *Status {
	st := &Status{
		eb: eventbus.New(),
	}
	st.Refresh(svcs)

	return st
}

// Update updates the status of the given service
func (st *Status) Update(service string, status HealthStatus) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.UpdatedAt = time.Now()
	for i, svc := range st.Report {
		if svc.Name == service {
			st.Report[i].HealthStatus = status
			st.eb.Publish(string(UPDATE), svc)
			return
		}
	}
}

// Refresh refreshes the status
func (st *Status) Refresh(svcs []config.Service) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.UpdatedAt = time.Now()
	st.Report = []ServiceHealth{}

	for _, svc := range svcs {
		st.Report = append(st.Report, ServiceHealth{
			Service:      svc,
			HealthStatus: HealthUnknown,
		})
	}
}

// Of method returns ServiceHealth of the service name
// given in the parameter
func (st *Status) Of(service string) ServiceHealth {
	for _, svc := range st.Report {
		if svc.Name == service {
			return svc
		}
	}

	return ServiceHealth{}
}

// Copy returns the value of the current status object
func (st *Status) Copy() *Status {
	st.mu.Lock()
	defer st.mu.Unlock()

	new := New(nil)
	new.Report = append(new.Report, st.Report...)

	return new
}

// Watch takes in a event name that can occur in the status object
// and a callback which will be invoked whenever the given event occurs
func (st *Status) Watch(ev Event, cb WatchCallback) *Watcher {
	id := st.eb.Subscribe(string(ev), func(data ...interface{}) {
		if len(data) != 1 {
			logrus.Errorf("malformed data received for event %s", ev)
			return
		}

		svc, ok := data[0].(ServiceHealth)
		if !ok {
			logrus.Errorf("malformed data received for event %s, expected ServiceHealth", ev)
			return
		}

		cb(svc)
	})

	return &Watcher{
		ID:    id,
		Topic: ev,
		eb:    st.eb,
	}
}
