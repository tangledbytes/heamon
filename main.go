package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

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

var (
	uiDirectory = filepath.Join(".", "ui", "build")
	version     = "dev"
	commit      = "none"
	date        = "NA"
)

func main() {
	printInfo(version, commit, date)

	// Setup the rendering engine as there are some overrides
	// that heamon offers react based frontend
	engine := html.New(uiDirectory, ".html")
	app := fiber.New(fiber.Config{Views: engine})

	// Setup monitoring on top of the http endpoints
	mon := monitor.New(models.Config{Interval: 1})
	handlers := handlers.NewHandlers(mon.Config, mon.Status)

	routes.NewRoutes(app, handlers)

	go gracefulShutdown(app)

	if err := app.Listen(":" + utils.GetEnv("PORT", defaultPort)); err != nil {
		logrus.Error(err)
	}
}

func gracefulShutdown(app *fiber.App) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	logrus.Info("Shutting down Heamon 👋")

	if err := app.Shutdown(); err != nil {
		logrus.Error("error occured while shutting down the Heamon:", err)
	}
}

func printInfo(version, commit, date string) {
	fmt.Println("🔥 Heamon version:", version)
	fmt.Println("🛠️ Commit:", commit)
	fmt.Println("📅 Release Date:", date)
}
