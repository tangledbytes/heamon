package handlers

import "github.com/gofiber/fiber/v2"

// GetStatus returns the status
func (h *Handler) GetStatus(c *fiber.Ctx) error {
	return c.JSON(h.status.GetStatus())
}
