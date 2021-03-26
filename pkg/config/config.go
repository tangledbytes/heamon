package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/utkarsh-pro/heamon/models"
)

// Setup will setup the config reader
func Setup() (cfg models.Config) {
	// Set defaults
	viper.SetDefault("HEAMON_USER", "admin")
	viper.SetDefault("HEAMON_PASS", "pl,pl,")
	viper.SetDefault("PORT", "5000")

	// Register alias to support mapping of env with config
	viper.RegisterAlias("HEAMON_USER", "authentication.username")
	viper.RegisterAlias("HEAMON_PASS", "authentication.password")

	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/heamon/")
	viper.AddConfigPath("$HOME/.heamon")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			logrus.Fatal("failed to read config:", err)
		}
	}

	viper.Unmarshal(&cfg)

	return
}
