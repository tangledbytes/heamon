package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/utkarsh-pro/heamon/pkg/utils"
)

// Setup sets up some default middlewares
func Setup(app *fiber.App) {
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(CustomBasicAuth(CustomBasicAuthConfig{
		Routes: Routes{"/api/v1/config"},
		Config: basicauth.Config{
			Users: map[string]string{
				utils.GetEnv("HEAMON_USER", "admin"): utils.GetEnv("HEAMON_PASS", "pl,pl,"),
			},
		},
	}))
}
