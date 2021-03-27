package routes

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/spf13/viper"
	"github.com/utkarsh-pro/heamon/pkg/entrypoint/rest/handlers"
)

// NewRoutes registers the REST endpoints to the fiber app
func NewRoutes(app *fiber.App, handlers *handlers.Handler, fsystem http.FileSystem) {
	apiV1Routes(app, handlers)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": viper.GetString("TITLE"),
		})
	})
	app.Use("/", filesystem.New(filesystem.Config{
		Root:   fsystem,
		Browse: true,
	}))
}

func apiV1Routes(app *fiber.App, handlers *handlers.Handler) {
	api := app.Group("/api/v1")

	api.Put("config", handlers.RegisterNewConfig)
	api.Get("config", handlers.GetConfig)
	api.Get("status", handlers.GetStatus)
}
