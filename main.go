package main

import (
	"encoding/json"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/utkarsh-pro/heamon/handlers"
	"github.com/utkarsh-pro/heamon/models"
	"github.com/utkarsh-pro/heamon/pkg/monitor"
	"github.com/utkarsh-pro/heamon/routes"
)

const defaultPort = "5000"

func main() {
	app := fiber.New()

	config := getBasicConfig()
	mon := monitor.New(config)
	handlers := handlers.NewHandlers(mon.Config, mon.Status)

	routes.NewRoutes(app, handlers)

	if err := app.Listen(":" + getPort()); err != nil {
		logrus.Error(err)
	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return defaultPort
	}

	return port
}

func getBasicConfig() models.Config {
	var cfg models.Config
	file, err := os.Open("./config.json")
	if err != nil {
		logrus.Error("failed to open config file", err.Error())
		return cfg
	}

	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		logrus.Error("failed to parse file data into config", err.Error())
	}

	return cfg
}
