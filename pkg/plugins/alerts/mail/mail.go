package mail

import (
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/utkarsh-pro/heamon/pkg/store"
	"github.com/utkarsh-pro/heamon/pkg/store/status"
)

// Mail struct
type Mail struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
	To       []string
	Duration float64

	status    store.Status
	watcher   *status.Watcher
	alertedOn *time.Time
}

// New returns a pointer to the instance of Mail Alert Struct
func New(host, port, password, username, from string, to []string, duration float64, status store.Status) *Mail {
	return &Mail{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
		To:       to,
		Duration: duration,
		status:   status,
	}
}

// Start will start the listener on status updates and
// will send email notifications if needed
func (m *Mail) Start() {
	m.watcher = m.status.Watch(status.UPDATE, func(sh status.ServiceHealth) {
		if sh.HealthStatus == status.HealthFail {
			t := time.Now()

			if m.alertedOn == nil ||
				(m.alertedOn.Add(time.Duration(m.Duration)*time.Minute).Sub(t) <= 0) {
				logrus.Info("[MAIL PLUGIN]: Sending Alert")
				m.Send(
					fmt.Sprintf("Service %s is failing", sh.Name),
					fmt.Sprintf(
						`Service %s with host name %s and health check endpoint %s is failing`,
						sh.Name,
						sh.Host,
						sh.HealthCheckEndpoint,
					),
				)
				m.alertedOn = &t
			}
		}
	})
}

func (m *Mail) Terminate() {
	m.watcher.Close()
}

// Send sends the content
func (m *Mail) Send(subject, message string) error {
	// Authentication
	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)

	content := fmt.Sprintf(`To: %s
From: %s
Subject: %s

%s`, strings.Join(m.To, ","), m.From, subject, message)

	// Sending Email
	return smtp.SendMail(m.Host+":"+m.Port, auth, m.From, m.To, []byte(content))
}
