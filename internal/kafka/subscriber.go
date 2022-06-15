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

type SubscriberClient struct {
	Reader *kaf.Reader
}

type Subscriber interface {
	Sub(ctx context.Context, event chan<- *pb.Event)
	Close() error
}

func NewSubscriber(topic string, addr string) Subscriber {
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

	return &SubscriberClient{Reader: reader}
}

func (s *SubscriberClient) Sub(ctx context.Context, event chan<- *pb.Event) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := s.Reader.ReadMessage(ctx)
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

func (s *SubscriberClient) Close() error {
	return s.Reader.Close()
}
