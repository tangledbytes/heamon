package rest

import (
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/utkarsh-pro/heamon/pkg/entrypoint/rest/handlers"
	"github.com/utkarsh-pro/heamon/pkg/entrypoint/rest/middlewares"
	"github.com/utkarsh-pro/heamon/pkg/entrypoint/rest/routes"
	"github.com/utkarsh-pro/heamon/pkg/monitor"
	"github.com/utkarsh-pro/heamon/pkg/store"
)

func Run() error {
	uiDirectory := filepath.Join(".", "ui", "build")

	manager := store.NewManager()
	manager.InitializeStore()

	mon := monitor.New(manager)

	engine := html.New(uiDirectory, ".html")
	app := fiber.New(fiber.Config{
		Views:                 engine,
		DisableStartupMessage: true,
	})

	middlewares.Setup(app)

	routes.NewRoutes(app, handlers.NewHandlers(mon.Config, mon.Status))

	return nil
}
