package plugins

import (
	"github.com/sirupsen/logrus"
	"github.com/utkarsh-pro/heamon/models"
	"github.com/utkarsh-pro/heamon/pkg/eventbus"
	"github.com/utkarsh-pro/heamon/plugins/alerts"
)

func Setup(cfg *models.PluginsAlerts) {
	if cfg.Email != nil {
		mail := alerts.NewMail(
			cfg.Email.SMTP.Host,
			cfg.Email.SMTP.Port,
			cfg.Email.SMTP.Password,
			cfg.Email.SMTP.Username,
			cfg.Email.From,
			cfg.Email.To,
		)

		eventbus.Bus.Subscribe(mail.GetEvent(), func(data ...interface{}) {
			if len(data) != 2 {
				logrus.Error("[Mail Alerts]: invalid message length")
				return
			}

			fails, ok := data[0].(float64)
			if !ok {
				logrus.Error("[Mail Alerts]: invalid failure rate type")
				return
			}

			service, ok := data[1].(string)
			if !ok {
				logrus.Error("[Mail Alerts]: invalid service name type")
				return
			}

			mail.EventCallback(fails, service)
		})
	}
}
