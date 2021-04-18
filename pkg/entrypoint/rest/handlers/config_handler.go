package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/utkarsh-pro/heamon/pkg/store/config"
	"github.com/valyala/fasthttp"
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

func (h *Handler) WatchConfig(c *fiber.Ctx) error {
	ctx := c.Context()

	ctx.SetContentType("text/event-stream")
	ctx.Response.Header.Set("Cache-Control", "no-cache")
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.Response.Header.Set("Transfer-Encoding", "chunked")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Cache-Control")
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	ctx.SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		ch := make(chan *config.Config, 1)

		h.config.Watch(config.UPDATE, func(c *config.Config) {
			ch <- c
		})

		for data := range ch {
			byt, err := json.Marshal(data)
			if err != nil {
				logrus.Error("failed to marshal data:", err.Error())
				continue
			}

			fmt.Fprintf(w, "%s", byt)

			if err := w.Flush(); err != nil {
				logrus.Error("error from stream writer:", err.Error())
				return
			}
		}
	}))

	println("exiting")

	return nil
}
