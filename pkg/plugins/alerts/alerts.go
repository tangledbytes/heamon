package alerts

import (
	"github.com/utkarsh-pro/heamon/pkg/plugins/alerts/mail"
	"github.com/utkarsh-pro/heamon/pkg/store"
	"github.com/utkarsh-pro/heamon/pkg/store/config"
)

type Alerts struct {
	config store.Config
	status store.Status

	mail *mail.Mail
}

func New(config store.Config, status store.Status) *Alerts {
	return &Alerts{
		config: config,
		status: status,
	}
}

func (a *Alerts) Start() {
	cfg := a.config.Copy()

	a.setupMailAlert(cfg)

	a.setupWatchAlertPlugins()
}

func (a *Alerts) setupWatchAlertPlugins() {
	a.config.Watch(config.UPDATE, func(c *config.Config) {
		a.setupMailAlert(c)
	})
}

func (a *Alerts) setupMailAlert(c *config.Config) {
	if c.Plugins != nil && c.Plugins.Alert != nil {
		if c.Plugins.Alert.Email != nil {
			// Terminate if instance already exists
			if a.mail != nil {
				a.mail.Terminate()
			}

			// Create new instance
			emailCfg := c.Plugins.Alert.Email
			smtpCfg := emailCfg.SMTP
			a.mail = mail.New(
				smtpCfg.Host,
				smtpCfg.Port,
				smtpCfg.Password,
				smtpCfg.Username,
				emailCfg.From,
				emailCfg.To,
				emailCfg.Duration,
				a.status,
			)

			// Start new instance
			a.mail.Start()
		}
	}
}
