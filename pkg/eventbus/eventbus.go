package eventbus

import (
	"sync"

	"github.com/utkarsh-pro/heamon/pkg/utils"
)

// EventBus represents an event bus which can be used for asynchronous
// communication among different components of the application
//
// It has a pub-sub based interface
type EventBus struct {
	subscribers map[string]map[int64]DataChannel
	rm          sync.RWMutex
}

// DataChannel - channel for exchanging data among publishers
// and subscribers
type DataChannel chan []interface{}

// DataEventCB - callback function type whenever a subscriber receives
// an event of interest
type DataEventCB func(data ...interface{})

// New returns an instance of event bus
func New() *EventBus {
	return &EventBus{
		subscribers: make(map[string]map[int64]DataChannel),
	}
}

// Subscribe - subscribes to the given topic and will execute the
// given callback whenever a publisher publishes on that topic
func (eb *EventBus) Subscribe(topic string, cb DataEventCB) (sid int64) {
	eb.rm.Lock()
	defer eb.rm.Unlock()

	ch := make(DataChannel, 5)
	topicGrp := eb.subscribers[topic]
	if topicGrp == nil {
		eb.subscribers[topic] = make(map[int64]DataChannel)
	}

	sid = utils.GenerateID()
	eb.subscribers[topic][sid] = ch

	// Start listening on the channel
	go func() {
		for data := range ch {
			cb(data...)
		}
	}()

	return sid
}

// Unsubscribe unsubscribe from the given topic
//
// sid - subscriber id
func (eb *EventBus) Unsubscribe(topic string, sid int64) {
	eb.rm.Lock()
	defer eb.rm.Unlock()

	topicGrp := eb.subscribers[topic]
	if topicGrp == nil {
		return
	}

	delete(eb.subscribers[topic], sid)
}

// Publish publishes data on the given topic
//
// Call is non-blocking
func (eb *EventBus) Publish(topic string, data ...interface{}) {
	go func() {
		eb.rm.Lock()
		defer eb.rm.Unlock()

		for _, ch := range eb.subscribers[topic] {
			ch <- data
		}
	}()
}

// Bus - eventbus singleton
var Bus = New()
