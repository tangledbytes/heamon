package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/utkarsh-pro/heamon/models"
)

// NewRoutes registers the REST endpoints to the fiber app
func NewRoutes(app *fiber.App, handlers models.Handler) {
	app.Put("/api/v1/config/monitor", handlers.RegisterNewConfig)
	app.Get("/api/v1/config", handlers.GetConfig)
	app.Get("/api/v1/status", handlers.GetStatus)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": viper.GetString("TITLE"),
		})
	})
	app.Static("/", "./ui/build")
}
