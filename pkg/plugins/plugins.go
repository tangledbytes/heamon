package plugins

import (
	"github.com/utkarsh-pro/heamon/pkg/plugins/alerts"
	"github.com/utkarsh-pro/heamon/pkg/store"
)

func Setup(config store.Config, status store.Status) {
	alerts.New(config, status).Start()
}
