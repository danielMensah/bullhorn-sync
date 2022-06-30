package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	kaf "github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

type KafkaPublisherClient struct {
	Conn *kaf.Conn
}

func NewKafkaPublisher(ctx context.Context, addr string) (Publisher, error) {
	conn, err := kaf.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return nil, fmt.Errorf("failed to set write deadline: %w", err)
	}

	// TODO: implement kafka dialer retry logic

	return &KafkaPublisherClient{Conn: conn}, nil
}

func (c *KafkaPublisherClient) Publish(ctx context.Context, events <-chan *EventWrapper, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-events:
			if !ok {
				return
			}

			marshalledData, err := json.Marshal(event.Data)
			if err != nil {
				log.WithError(err).Error("marshalling event")
				return
			}

			_, err = c.Conn.WriteMessages(kaf.Message{
				Value: marshalledData,
				Topic: event.Topic,
			})
			if err != nil {
				log.WithError(err).Error("writing message")
				return
			}
		}
	}
}

func (c *KafkaPublisherClient) Close() error {
	return c.Conn.Close()
}
