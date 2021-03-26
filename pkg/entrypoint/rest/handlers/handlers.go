package handlers

import (
	"github.com/utkarsh-pro/heamon/pkg/store"
)

type Handler struct {
	config store.Config
	status store.Status
}

// NewHandlers takes in a configuration and returns a
// pointer to an instance of handlers
func NewHandlers(config store.Config, initialStatus store.Status) *Handler {
	return &Handler{config: config, status: initialStatus}
}
