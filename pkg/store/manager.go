package store

import (
	"github.com/sirupsen/logrus"
	"github.com/utkarsh-pro/heamon/pkg/store/config"
	"github.com/utkarsh-pro/heamon/pkg/store/status"
)

type Manager struct {
	config *config.Config
	status *status.Status
}

// NewManager returns a pointer to an instance of Manager
func NewManager() *Manager {
	return &Manager{}
}

// Initialize initializes the entire store
// and ensures data consistency among the internal
// data collections
func (m *Manager) InitializeStore() {
	cfg := config.New()

	// Register config update hook
	cfg.Hook().Update.Register(func() {
		// Ensure that status and config object are ALWAYS
		// in sync
		m.status.Refresh(m.config.Monitor.Services)
	})

	initConfig(cfg)

	if err := cfg.Validate(); err != nil {
		logrus.Fatal("config validation failed:\n", err)
	}

	m.config = cfg
	m.status = status.New(cfg.Monitor.Services)
}

// Config returns an instance to the config object
func (m *Manager) Config() Config {
	return m.config
}

// Status returns an instance to the status object
func (m *Manager) Status() Status {
	return m.status
}
