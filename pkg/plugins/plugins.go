package plugins

import (
	"github.com/sirupsen/logrus"
	"github.com/utkarsh-pro/heamon/pkg/eventbus"
	"github.com/utkarsh-pro/heamon/pkg/plugins/alerts"
	"github.com/utkarsh-pro/heamon/pkg/store/config"
)

func Setup(cfg *config.Plugins) {
	if cfg.Alert.Email != nil {
		mail := alerts.NewMail(
			cfg.Alert.Email.SMTP.Host,
			cfg.Alert.Email.SMTP.Port,
			cfg.Alert.Email.SMTP.Password,
			cfg.Alert.Email.SMTP.Username,
			cfg.Alert.Email.From,
			cfg.Alert.Email.To,
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
