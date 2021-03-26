package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/viper"
)

// Setup sets up some default middlewares
func Setup(app *fiber.App) {
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(CustomBasicAuth(CustomBasicAuthConfig{
		Routes: Routes{"/api/v1/config/monitor"},
		Config: basicauth.Config{
			Users: map[string]string{
				viper.GetString("HEAMON_USER"): viper.GetString("HEAMON_PASS"),
			},
		},
	}))
}
