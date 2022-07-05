package broker

import (
	"context"
	"sync"
)

// EventWrapper wraps an event with a topic
type EventWrapper struct {
	Topic string
	Event interface{} `json:"event"`
}

// Consumer is an interface that defines the methods that a consumer must implement
type Consumer interface {
	// Consume consumes events from a topic
	Consume(ctx context.Context, topic string, event chan<- *EventWrapper)
	// Close closes the consumer
	Close() error
}

// Producer is an interface that defines the methods that a producer must implement
type Producer interface {
	// MonitorEvents monitors events from the producer
	MonitorEvents()
	// Produce sends events to the producer
	Produce(ctx context.Context, events <-chan *EventWrapper, wg *sync.WaitGroup)
	// Close closes the producer
	Close()
}
