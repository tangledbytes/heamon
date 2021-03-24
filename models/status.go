package models

import (
	"time"
)

// HealthStatus is custom type for indicating a service's health
type HealthStatus string

const (
	HealthUnknown      HealthStatus = "UNKNOWN"
	HealthOK           HealthStatus = "OK"
	HealthFail         HealthStatus = "FAIL"
	Healthdeteriorated HealthStatus = "DETERIORATED"
)

// Status stores the status of the services
type Status struct {
	UpdatedAt time.Time       `json:"updated_at,omitempty"`
	Report    []ServiceHealth `json:"report,omitempty"`
}

// ServiceHealth captures the metadata about the service
// along with its status
type ServiceHealth struct {
	Service      `json:"service,omitempty"`
	HealthStatus HealthStatus `json:"health_status,omitempty"`
}

// NewStatus returns a new instance of the status
func NewStatus(svcs []Service) *Status {
	st := &Status{}
	st.Refresh(svcs)

	return st
}

// UpdateStatus updates the status of the given service
func (st *Status) UpdateStatus(service string, status HealthStatus) {
	st.UpdatedAt = time.Now()
	for i, svc := range st.Report {
		if svc.Name == service {
			st.Report[i].HealthStatus = status
			return
		}
	}
}

// Refresh refreshes the status
func (st *Status) Refresh(svcs []Service) {
	st.UpdatedAt = time.Now()
	st.Report = []ServiceHealth{}

	for _, svc := range svcs {
		st.Report = append(st.Report, ServiceHealth{
			Service:      svc,
			HealthStatus: HealthUnknown,
		})
	}
}

// GetStatus returns the value of the current status object
func (st *Status) GetStatus() Status {
	return *st
}
