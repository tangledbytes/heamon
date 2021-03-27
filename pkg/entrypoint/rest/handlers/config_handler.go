package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// RegisterNewConfig updates the configuration
func (h *Handler) RegisterNewConfig(c *fiber.Ctx) error {
	if err := h.config.Update(c.Body()); err != nil {
		return c.Status(http.StatusBadRequest).JSON(GenericMessageResponse{err.Error()})
	}

	return c.Status(http.StatusOK).JSON(GenericMessageResponse{"successfully updated"})
}

// GetConfig returns the current config
func (h *Handler) GetConfig(c *fiber.Ctx) error {
	return c.JSON(h.config.Copy())
}

// PatchConfig patches the configuration
func (h *Handler) PatchConfig(c *fiber.Ctx) error {
	if err := h.config.Merge(c.Body()); err != nil {
		return c.Status(http.StatusBadRequest).JSON(GenericMessageResponse{err.Error()})
	}

	return c.Status(http.StatusOK).JSON(GenericMessageResponse{"successfully updated"})
}
