package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

// Routes type is an alias for string string type
// and adds some additional functions to it
type Routes []string

// IsIn returns true if the given route is in the Routes
func (rs Routes) IsIn(route string) bool {
	for _, r := range rs {
		if string(r) == route {
			return true
		}
	}

	return false
}

// Config defines the config for middleware.
type CustomBasicAuthConfig struct {
	basicauth.Config
	Routes Routes
}

// CustomBasicAuth is a HOF, it will take in a slice
// of routes and will return a middleware
func CustomBasicAuth(config CustomBasicAuthConfig) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ba := basicauth.New(config.Config)
		if config.Routes.IsIn(c.Path()) {
			return ba(c)
		}

		return c.Next()
	}
}
