package config

import "github.com/utkarsh-pro/heamon/pkg/eventbus"

// Event is type for the events that Config struct's
// Watch method can consume
type Event string

const (
	// UPDATE is fired whenever Update method
	// is called on the config object
	UPDATE Event = "heamon.config.update"
)

// WatchCallback is the type for callback functions
// that Watch method of Config struct can consume
type WatchCallback func()

// Watcher struct manages something -_-
type Watcher struct {
	ID    int64
	Topic Event
	eb    *eventbus.EventBus
}

// Close will "unsubscribe" the watch
//
// It is RECOMMENDED to invoke the function
// if the consumer is no longer interested in
// the event
func (w *Watcher) Close() {
	w.eb.Unsubscribe(string(w.Topic), w.ID)
}
