package rest

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/utkarsh-pro/heamon/pkg/entrypoint/rest/handlers"
	"github.com/utkarsh-pro/heamon/pkg/entrypoint/rest/middlewares"
	"github.com/utkarsh-pro/heamon/pkg/entrypoint/rest/routes"
	"github.com/utkarsh-pro/heamon/pkg/monitor"
	"github.com/utkarsh-pro/heamon/pkg/plugins"
	"github.com/utkarsh-pro/heamon/pkg/store"
)

//go:embed ui
var embededFiles embed.FS

func Run() error {
	manager := store.NewManager()
	manager.InitializeStore()

	mon := monitor.New(manager.Config(), manager.Status())

	plugins.Setup(manager.Config(), manager.Status())

	fsystem := getFileSystem()
	engine := html.NewFileSystem(fsystem, ".html")
	app := fiber.New(fiber.Config{
		Views:                 engine,
		DisableStartupMessage: true,
		ReadTimeout:           30 * time.Second,
	})

	middlewares.Setup(app)

	routes.NewRoutes(app, handlers.NewHandlers(mon.Config, mon.Status), fsystem)

	go func() {
		if err := app.Listen(":" + viper.GetString("PORT")); err != nil {
			logrus.Error(err)
		}
	}()

	fmt.Println("🚀 Heamon started on PORT: ", viper.GetString("PORT"))

	// Handle graceful shutdown
	gracefulShutdown(app, mon.Stop)

	return nil
}

func getFileSystem() http.FileSystem {
	fsys, err := fs.Sub(embededFiles, "ui")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

func gracefulShutdown(app *fiber.App, cleanup func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	logrus.Info("Shutting down Heamon 👋")

	cleanup()

	if err := app.Shutdown(); err != nil {
		logrus.Error("error occured while shutting down the Heamon:", err)
	}
}
