package store

import (
	"github.com/utkarsh-pro/heamon/pkg/store/config"
)

// Config is an interface for Config object
//
// It allows to give limited control over the Config
// object while still able to perform write operations
type Config interface {
	Copy() *config.Config
	Update([]byte) error
	Hook() *config.Hook
	Watch(config.Event, config.WatchCallback) *config.Watcher
	Merge([]byte) error
}
