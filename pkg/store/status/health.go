package status

// HealthStatus is custom type for indicating a service's health
type HealthStatus string

const (
	HealthUnknown  HealthStatus = "UNKNOWN"
	HealthOK       HealthStatus = "OK"
	HealthFail     HealthStatus = "FAIL"
	HealthDegraded HealthStatus = "DEGRADED"
)
