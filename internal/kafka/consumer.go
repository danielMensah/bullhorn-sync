package kafka

import (
	"context"
	"net"
	"time"

	pb "github.com/danielMensah/bullhorn-sync-poc/internal/proto"
	"github.com/golang/protobuf/proto"
	kaf "github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

type ConsumerClient struct {
	reader *kaf.Reader
}

type Consumer interface {
	Consume(ctx context.Context, event chan<- *pb.Event)
	Close() error
}

func NewConsumer(topic string, addr string) Consumer {
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

	return &ConsumerClient{reader: reader}
}

func (s *ConsumerClient) Consume(ctx context.Context, event chan<- *pb.Event) {
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

			e := &pb.Event{}
			err = proto.Unmarshal(msg.Value, e)
			if err != nil {
				log.WithError(err).Error("failed to unmarshal message")
			}

			event <- e
		}
	}
}

func (s *ConsumerClient) Close() error {
	return s.reader.Close()
}
