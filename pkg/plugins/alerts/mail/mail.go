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

	status     store.Status
	watcher    *status.Watcher
	alertedOn  *time.Time
	alertedFor status.HealthStatus
}

// New returns a pointer to the instance of Mail Alert Struct
func New(host, port, password, username, from string, to []string, duration float64, status store.Status) *Mail {
	return &Mail{
		Host:       host,
		Port:       port,
		Username:   username,
		Password:   password,
		From:       from,
		To:         to,
		Duration:   duration,
		status:     status,
		alertedFor: "",
	}
}

// Start will start the listener on status updates and
// will send email notifications if needed
func (m *Mail) Start() {
	m.watcher = m.status.Watch(status.UPDATE, func(sh status.ServiceHealth) {
		const msg = `
Service %s with host name %s and health checkpoint %s is in state "%s".

Current Heamon config:
	Service State "FAIL" over Failure Rate: %f
	Service State "DEGRADED" over Failure Rate: %f
`
		if sh.HealthStatus == status.HealthFail || sh.HealthStatus == status.HealthDegraded {
			t := time.Now()

			if m.alertedOn == nil ||
				(m.alertedOn.Add(time.Duration(m.Duration)*time.Minute).Sub(t) <= 0) ||
				m.alertedFor != sh.HealthStatus {
				logrus.Info("[MAIL PLUGIN]: Sending Alert")

				if err := m.Send(
					fmt.Sprintf("Alert For Service %s", sh.Name),
					fmt.Sprintf(msg,
						sh.Name,
						sh.Host,
						sh.HealthCheckEndpoint,
						sh.HealthStatus,
						sh.Failure,
						sh.Degraded,
					),
				); err != nil {
					logrus.Error("[MAIL]:", err)
				}

				m.alertedOn = &t
				m.alertedFor = sh.HealthStatus
			}
		}
	})
}

// Terminate will gracefully shutdown the mail plugin
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
