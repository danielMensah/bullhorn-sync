package broker

import (
	"context"
	"encoding/json"
	"net"
	"time"

	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

type KafkaConsumerClient struct {
	reader KafkaReaderService
}

// KafkaReaderService is for mainly testing purposes
type KafkaReaderService interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	Close() error
}

func NewKafkaConsumer(topic string, addr string) Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{addr},
		Topic:   topic,
		Dialer: &kafka.Dialer{
			DialFunc: func(ctx context.Context, network string, address string) (net.Conn, error) {
				return kafka.DialContext(ctx, "tcp", addr)
			},
			Timeout: 10 * time.Second,
		},
	})

	return &KafkaConsumerClient{reader: reader}
}

func (s *KafkaConsumerClient) Consume(ctx context.Context, event chan<- *EventWrapper) {
	for {
		select {
		case <-ctx.Done():
			close(event)
			return
		default:
			msg, err := s.reader.ReadMessage(ctx)
			if err != nil {
				log.WithError(err).Error("failed to read message")
				return
			}

			e := &EventWrapper{}
			err = json.Unmarshal(msg.Value, e)
			if err != nil {
				log.WithError(err).Error("failed to unmarshal message")
				return
			}

			event <- e
		}
	}
}

func (s *KafkaConsumerClient) Close() error {
	return s.reader.Close()
}
