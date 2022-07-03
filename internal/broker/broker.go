package broker

import (
	"context"
	"sync"
)

type EventWrapper struct {
	Topic string
	Event interface{} `json:"event"`
}

type Consumer interface {
	Consume(ctx context.Context, topic string, event chan<- *EventWrapper)
}

type Producer interface {
	MonitorEvents(wg *sync.WaitGroup)
	Produce(ctx context.Context, events <-chan *EventWrapper, wg *sync.WaitGroup)
}
