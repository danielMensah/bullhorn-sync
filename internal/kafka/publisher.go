package kafka

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/danielMensah/bullhorn-sync-poc/internal/bullhorn"
	pb "github.com/danielMensah/bullhorn-sync-poc/internal/proto"
	kaf "github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type PublisherClient struct {
	Conn *kaf.Conn
}

type Publisher interface {
	Publish(ctx context.Context, event <-chan *bullhorn.Entity, wg *sync.WaitGroup)
	Close() error
}

func NewPublisher(ctx context.Context, addr string) (Publisher, error) {
	conn, err := kaf.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return nil, fmt.Errorf("failed to set write deadline: %w", err)
	}

	// TODO: implement kafka dialer retry logic

	return &PublisherClient{Conn: conn}, nil
}

func (c *PublisherClient) Publish(ctx context.Context, entityEvent <-chan *bullhorn.Entity, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case entity, ok := <-entityEvent:
			if !ok {
				return
			}

			p := &pb.Entity{
				Id:             entity.Id,
				Name:           entity.Name,
				Changes:        entity.Changes,
				EventTimestamp: entity.Timestamp,
			}

			data, err := proto.Marshal(p)
			if err != nil {
				log.WithError(err).Error("marshalling proto event")
				return
			}

			_, err = c.Conn.WriteMessages(kaf.Message{
				Value: data,
				Topic: fmt.Sprintf("event_%s", entity.Name),
			})
			if err != nil {
				log.WithError(err).Error("writing message")
				return
			}
		}
	}
}

func (c *PublisherClient) Close() error {
	return c.Conn.Close()
}
