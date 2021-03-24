package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/sirupsen/logrus"
	"github.com/utkarsh-pro/heamon/handlers"
	"github.com/utkarsh-pro/heamon/models"
	"github.com/utkarsh-pro/heamon/pkg/monitor"
	"github.com/utkarsh-pro/heamon/pkg/utils"
	"github.com/utkarsh-pro/heamon/routes"
)

const defaultPort = "5000"

func main() {
	engine := html.New("./ui/build", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	config := getBasicConfig()
	mon := monitor.New(config)
	handlers := handlers.NewHandlers(mon.Config, mon.Status)

	routes.NewRoutes(app, handlers)

	if err := app.Listen(":" + utils.GetEnv("PORT", defaultPort)); err != nil {
		logrus.Error(err)
	}
}

func getBasicConfig() models.Config {
	return models.Config{
		Interval: 1,
	}
}
