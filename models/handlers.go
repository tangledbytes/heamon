package models

import "github.com/gofiber/fiber/v2"

// Handler is the interface for the handler
type Handler interface {
	RegisterNewConfig(c *fiber.Ctx) error
	GetConfig(c *fiber.Ctx) error
	GetStatus(c *fiber.Ctx) error
}
