package config

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/utkarsh-pro/heamon/pkg/eventbus"
	"github.com/utkarsh-pro/heamon/pkg/hook"
)

// Config stores the configuration
type Config struct {
	Title          string         `json:"title,omitempty"`
	Port           string         `json:"port,omitempty"`
	Authentication Authentication `json:"authentication,omitempty"`

	Monitor Monitor  `json:"monitor,omitempty"`
	Plugins *Plugins `json:"plugins,omitempty"`

	hook *Hook `json:"-"`
	mu   sync.Mutex
	eb   *eventbus.EventBus
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

// New returns a pointer to an instance
// of the Config struct
//
// This instance isn't initialized and should be
// initialized using Update method
func New() *Config {
	return &Config{
		hook: &Hook{
			Update: hook.New(),
		},
	}
}

// Validate returns an error if the configuration is invalid
func (cfg *Config) Validate() error {
	return nil
}

// Update updates the config
func (cfg *Config) Update(configbyt []byte) error {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	var tempCfg Config
	if err := json.Unmarshal(configbyt, &tempCfg); err != nil {
		return fmt.Errorf("invalid format of the configuration")
	}

	if err := tempCfg.Validate(); err != nil {
		return err
	}

	if err := json.Unmarshal(configbyt, cfg); err != nil {
		return fmt.Errorf("invalid format of the configuration")
	}

	// Fire the event
	cfg.eb.Publish(string(UPDATE))

	// Execute the update hook
	cfg.Hook().Update.Execute()

	return nil
}

// Hook returns a pointer to the hook
//
// This can be used to register update hooks
func (cfg *Config) Hook() *Hook {
	return cfg.hook
}

// Copy will return a copy of the config object
//
// Even though this will return a pointer
// to the Config object but THIS WILL STILL
// BE A COPY of the original config
//
// The copy also doesn't have a valid Hook in place
func (cfg *Config) Copy() *Config {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	plugins := *cfg.Plugins

	ccfg := &Config{
		Title:          cfg.Title,
		Port:           cfg.Port,
		Authentication: cfg.Authentication,

		Monitor: cfg.Monitor,
		Plugins: &plugins,
	}

	return ccfg
}

func (cfg *Config) Watch(ev Event, cb WatchCallback) *Watcher {
	id := cfg.eb.Subscribe(string(ev), func(data ...interface{}) {
		cb()
	})

	return &Watcher{
		ID:    id,
		Topic: ev,
		eb:    cfg.eb,
	}
}
