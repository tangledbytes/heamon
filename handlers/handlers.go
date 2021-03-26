package handlers

import (
	"github.com/utkarsh-pro/heamon/models"
)

type Status interface {
	Refresh(svcs []models.Service)
	GetStatus() models.Status
}

type Config interface {
	UpdateConfig([]byte) error
}

type Handler struct {
	config Config
	status Status
}

// NewHandlers takes in a configuration and returns a
// pointer to an instance of handlers
func NewHandlers(config Config, initialStatus Status) *Handler {
	return &Handler{config: config, status: initialStatus}
}
