package kafka

import (
	"context"
	"fmt"
	"time"

	pb "github.com/danielMensah/bullhorn-sync-poc/internal/proto"
	"github.com/golang/protobuf/proto"
	kaf "github.com/segmentio/kafka-go"
)

type Client struct {
	Conn *kaf.Conn
}

type Kafka interface {
	Close() error
	Pub(events []*pb.Event) error
}

func NewMessenger(ctx context.Context, addr string) (Kafka, error) {
	conn, err := kaf.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return nil, fmt.Errorf("failed to set write deadline: %w", err)
	}

	// TODO: implement kafka dialer retry logic

	return &Client{Conn: conn}, nil
}

func (c *Client) Pub(events []*pb.Event) error {
	messages := make([]kaf.Message, len(events))
	for _, event := range events {
		data, err := proto.Marshal(event)
		if err != nil {
			return fmt.Errorf("failed to marshal event: %w", err)
		}

		messages = append(messages, kaf.Message{
			Value: data,
			Topic: fmt.Sprintf("event_%s", event.EntityEventType),
		})
	}

	_, err := c.Conn.WriteMessages(messages...)
	if err != nil {
		return fmt.Errorf("failed to write messages: %w", err)
	}

	return nil
}

func (c *Client) Close() error {
	return c.Conn.Close()
}
