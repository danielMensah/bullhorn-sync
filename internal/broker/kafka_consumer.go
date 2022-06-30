package broker

import (
	"context"
	"encoding/json"
	"net"
	"time"

	kaf "github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

type KafkaConsumerClient struct {
	reader *kaf.Reader
}

func NewKafkaConsumer(topic string, addr string) Consumer {
	reader := kaf.NewReader(kaf.ReaderConfig{
		Brokers: []string{addr},
		Topic:   topic,
		Dialer: &kaf.Dialer{
			DialFunc: func(ctx context.Context, network string, address string) (net.Conn, error) {
				return kaf.DialContext(ctx, "tcp", addr)
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
			}

			e := &EventWrapper{}
			err = json.Unmarshal(msg.Value, e)
			if err != nil {
				log.WithError(err).Error("failed to unmarshal message")
			}

			event <- e
		}
	}
}

func (s *KafkaConsumerClient) Close() error {
	return s.reader.Close()
}
