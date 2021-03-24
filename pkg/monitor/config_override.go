package monitor

import (
	"github.com/utkarsh-pro/heamon/models"
	"github.com/utkarsh-pro/heamon/pkg/eventbus"
)

type Config struct {
	*models.Config
}

func (cfg *Config) UpdateConfig(configByt []byte) error {
	// Update the config
	if err := cfg.Config.UpdateConfig(configByt); err != nil {
		return err
	}

	// Publish the event
	eventbus.Bus.Publish(ConfigUpdate)

	return nil
}
