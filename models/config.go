package models

import (
	"encoding/json"
	"fmt"
)

type Config struct {
	Title          string         `json:"title,omitempty"`
	Port           string         `json:"port,omitempty"`
	Authentication Authentication `json:"authentication,omitempty"`

	Monitor Monitor  `json:"monitor,omitempty"`
	Plugins *Plugins `json:"plugins,omitempty"`
}

type Authentication struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type Plugins struct {
	Alert *PluginsAlerts `json:"alert,omitempty"`
}

type PluginsAlerts struct {
	Email *AlertEmail `json:"email,omitempty"`
}

type AlertEmail struct {
	SMTP AlertEmailSMTP `json:"smtp,omitempty"`
	From string         `json:"from,omitempty"`
	To   []string       `json:"to,omitempty"`
}

type AlertEmailSMTP struct {
	Host     string `json:"host,omitempty"`
	Port     string `json:"port,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type Monitor struct {
	Interval    int       `json:"interval,omitempty"`
	AverageOver int       `json:"average_over,omitempty"`
	Services    []Service `json:"services,omitempty"`
}

type Service struct {
	Name                string  `json:"name,omitempty"`
	Host                string  `json:"host,omitempty"`
	Interval            int     `json:"interval,omitempty"`
	HealthCheckEndpoint string  `json:"health_check_endpoint,omitempty"`
	Tolerance           float64 `json:"tolerance,omitempty"`
}

// IsValid returns an error if the configuration is invalid
func (cfg *Config) IsValid() error {
	return nil
}

// UpdateConfig updates the config
func (cfg *Config) UpdateConfig(configbyt []byte) error {
	var tempCfg Config
	if err := json.Unmarshal(configbyt, &tempCfg); err != nil {
		return fmt.Errorf("invalid format of the configuration")
	}

	if err := tempCfg.IsValid(); err != nil {
		return err
	}

	if err := json.Unmarshal(configbyt, cfg); err != nil {
		return fmt.Errorf("invalid format of the configuration")
	}
	return nil
}
