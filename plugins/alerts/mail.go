package alerts

import "net/smtp"

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
func (m *Mail) Send(content []byte) error {
	// Authentication
	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)

	// Sending Email
	return smtp.SendMail(m.Host+":"+m.Port, auth, m.From, m.To, content)
}
