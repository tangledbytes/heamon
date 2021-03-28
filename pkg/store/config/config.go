package config

import (
	"encoding/json"
	"fmt"
	"sync"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/utkarsh-pro/heamon/pkg/eventbus"
	"github.com/utkarsh-pro/heamon/pkg/hook"
)

// Config stores the configuration300000
type Config struct {
	Title          string         `json:"title,omitempty" validate:"required"`
	Port           string         `json:"port,omitempty" validate:"required,number"`
	Authentication Authentication `json:"authentication,omitempty" validate:"required"`

	Monitor Monitor  `json:"monitor,omitempty" validate:"required"`
	Plugins *Plugins `json:"plugins,omitempty"`

	hook *Hook              `json:"-"`
	mu   sync.Mutex         `json:"-"`
	eb   *eventbus.EventBus `json:"-"`
}

type Authentication struct {
	Username string `json:"username,omitempty" validate:"required,min=5"`
	Password string `json:"password,omitempty" validate:"required,min=8"`
}

type Plugins struct {
	Alert *PluginsAlerts `json:"alert,omitempty" validate:"dive"`
}

type PluginsAlerts struct {
	Email *AlertEmail `json:"email,omitempty" validate:"dive"`
}

type AlertEmail struct {
	SMTP     AlertEmailSMTP `json:"smtp,omitempty" validate:"required,dive"`
	From     string         `json:"from,omitempty" validate:"required,startswith=<,endswith=>"`
	To       []string       `json:"to,omitempty" validate:"required,dive,email"`
	Duration float64        `json:"duration,omitempty" validate:"gte=0"`
}

type AlertEmailSMTP struct {
	Host     string `json:"host,omitempty" validate:"required"`
	Port     string `json:"port,omitempty" validate:"required"`
	Username string `json:"username,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

type Monitor struct {
	Interval float64   `json:"interval,omitempty" validate:"gte=1"`
	Services []Service `json:"services,omitempty" validate:"unique=Name,dive"`
}

type Service struct {
	Name                string  `json:"name,omitempty" validate:"required"`
	Host                string  `json:"host,omitempty" validate:"required,hostname|hostname_port"`
	Interval            float64 `json:"interval,omitempty" validate:"-"`
	HealthCheckEndpoint string  `json:"health_check_endpoint,omitempty"`
	Failure             float64 `json:"failure,omitempty" validate:"gte=0,lte=100"`
	Degraded            float64 `json:"degraded,omitempty" validate:"gte=0,lte=100"`
	InitialDownTime     float64 `json:"initial_down_time,omitempty" validate:"gte=0"`
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
		eb: eventbus.New(),
	}
}

// Validate returns an error if the configuration is invalid
func (cfg *Config) Validate() error {
	validate := validator.New()

	return validate.Struct(cfg)
}

// Merge takes in patch config bytes and updates the config accordingly
func (cfg *Config) Merge(mergeByt []byte) error {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	// Marshal the config struct
	configByt, err := json.Marshal(cfg)
	if err != nil {
		logrus.Error("[Config Patch]:", err)
		return fmt.Errorf("failed to parse internal config: %s", err.Error())
	}

	// Merge the config
	final, err := jsonpatch.MergePatch(configByt, mergeByt)
	if err != nil {
		logrus.Error("[Config Patch]:", err)
		return fmt.Errorf("failed to merge the configuration: %s", err.Error())
	}

	return cfg.updateWithoutLock(final)
}

// Update updates the config
func (cfg *Config) Update(configbyt []byte) error {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	return cfg.updateWithoutLock(configbyt)
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

	ccfg := &Config{
		Title:          cfg.Title,
		Port:           cfg.Port,
		Authentication: cfg.Authentication,

		Monitor: cfg.Monitor,
	}

	if cfg.Plugins != nil {
		plugins := *cfg.Plugins
		ccfg.Plugins = &plugins
	}

	return ccfg
}

// Watch takes in a event name that can occur in the config object
// and a callback which will be invoked whenever the given event occurs
func (cfg *Config) Watch(ev Event, cb WatchCallback) *Watcher {
	id := cfg.eb.Subscribe(string(ev), func(data ...interface{}) {
		if len(data) != 1 {
			logrus.Errorf("malformed data received for event %s", ev)
			return
		}

		cfg, ok := data[0].(*Config)
		if !ok {
			logrus.Errorf("malformed data received for event %s, expected *Config", ev)
			return
		}

		cb(cfg)
	})

	return &Watcher{
		ID:    id,
		Topic: ev,
		eb:    cfg.eb,
	}
}

// updateWithoutLock is an internal function which attempts to update the internal
// configuration object WITHOUT acquiring the lock. This method EXPECTS that the caller
// would have already acquired lock on the object
func (cfg *Config) updateWithoutLock(configbyt []byte) error {
	var tempCfg Config
	if err := json.Unmarshal(configbyt, &tempCfg); err != nil {
		logrus.Error("[Config Update]:", err)
		return fmt.Errorf("invalid format of the configuration: %s", err.Error())
	}

	if err := tempCfg.Validate(); err != nil {
		return err
	}

	cfg.refresh()

	if err := json.Unmarshal(configbyt, cfg); err != nil {
		logrus.Error("[Config Update]:", err)
		return fmt.Errorf("invalid format of the configuration: %s", err.Error())
	}

	// Execute the update hook
	cfg.Hook().Update.Execute()

	// Fire the event THIS SHOULD HAPPEN ONLY
	// ONCE ALL OF THE HOOKS ARE COMPLETED
	cfg.eb.Publish(string(UPDATE), &tempCfg)

	return nil
}

// refresh clears the data from the config object
// so that fresh data can be entered into it
//
// This method IS NOT thread safe and hence should be invoked
// only by the caller which is acquiring lock on the config object
func (cfg *Config) refresh() {
	cfg.Title = ""
	cfg.Port = ""
	cfg.Authentication = Authentication{}
	cfg.Monitor = Monitor{}
	cfg.Plugins = nil
}
