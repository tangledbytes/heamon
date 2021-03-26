package config

import "github.com/utkarsh-pro/heamon/pkg/hook"

// Hook acts as an container for all the hooks associated
// with the config object
type Hook struct {
	Update *hook.Hook
}
