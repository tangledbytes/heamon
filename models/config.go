package models

import (
	"encoding/json"
	"fmt"
)

type Config struct {
	Interval int       `json:"interval,omitempty"`
	Services []Service `json:"services,omitempty"`
}

type Service struct {
	Name                string `json:"name,omitempty"`
	Host                string `json:"host,omitempty"`
	HealthCheckEndpoint string `json:"health_check_endpoint,omitempty"`
}

func (cfg *Config) GetConfig() Config {
	return *cfg
}

// IsValid returns an error if the configuration is invalid
func (cfg *Config) IsValid() error {
	// Check if the interval is greater than 0
	if cfg.Interval < 1 {
		return fmt.Errorf("interval cannot be less than 1 minute")
	}

	// Check for duplicate names
	names := map[string]struct{}{}
	for _, svc := range cfg.Services {
		_, ok := names[svc.Name]
		if ok {
			return fmt.Errorf("duplicate entries found for %s service", svc.Name)
		}

		names[svc.Name] = struct{}{}
	}

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
