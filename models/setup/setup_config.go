package setup

type Config struct {
	Title          string
	Port           string
	Authentication Authentication

	Plugins Plugins
}

type Authentication struct {
	Username string
	Password string
}

type Plugins struct {
	Alert PluginsAlerts
}

type PluginsAlerts struct {
	Email AlertEmail
}

type AlertEmail struct {
	Username string
	SMTP     AlertEmailSMTP
	To       []string
}

type AlertEmailSMTP struct {
	Host     string
	Port     string
	Password string
}
