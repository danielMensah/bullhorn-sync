package broker

import (
	"context"
	"sync"
)

type EventWrapper struct {
	Topic string
	Data  interface{}
}

type Consumer interface {
	Consume(ctx context.Context, event chan<- *EventWrapper)
	Close() error
}

type Publisher interface {
	Publish(ctx context.Context, events <-chan *EventWrapper, wg *sync.WaitGroup)
	Close() error
}
