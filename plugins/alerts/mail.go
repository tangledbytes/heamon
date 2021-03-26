package alerts

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/sirupsen/logrus"
)

// Mail struct
type Mail struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
	To       []string
}

// NewMail returns a pointer to the instance of Mail Alert Struct
func NewMail(host, port, password, username, from string, to []string) *Mail {
	return &Mail{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
		To:       to,
	}
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

func (m *Mail) GetEvent() string {
	return "heamon.plugin.alert"
}

func (m *Mail) EventCallback(fail float64, service string) {
	logrus.Info("sending alert mail!")
	if err := m.Send(service+" Alert!", fmt.Sprintf("%s has failed for over %f%% times!\n", service, fail)); err != nil {
		logrus.Error(err)
	}
}
