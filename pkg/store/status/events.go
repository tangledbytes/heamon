package status

import "github.com/utkarsh-pro/heamon/pkg/eventbus"

// Event is type for the events that Status struct's
// Watch method can consume
type Event string

const (
	// UPDATE is fired whenever Update method
	// is called on the struct object
	UPDATE Event = "heamon.status.update"
)

// WatchCallback is the type for callback functions
// that Watch method of Struct struct can consume
type WatchCallback func(HealthStatus)

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
