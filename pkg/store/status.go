package store

import (
	"github.com/utkarsh-pro/heamon/pkg/store/config"
	"github.com/utkarsh-pro/heamon/pkg/store/status"
)

// Status is an interface for status object
// allowing only limited interaction with the
// internal data in order to protects its integrity
type Status interface {
	Update(service string, status status.HealthStatus)
	Refresh(svcs []config.Service)
	Copy() *status.Status
	Watch(status.Event, status.WatchCallback) *status.Watcher
}
