package broker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
)

type kafkaConsumerClient struct {
	svc KafkaConsumerService
}

// KafkaConsumerService is an interface to mock the kafka consumer for testing
type KafkaConsumerService interface {
	ReadMessage(timeout time.Duration) (*kafka.Message, error)
	SubscribeTopics(topics []string, rebalanceCb kafka.RebalanceCb) (err error)
	Close() error
}

// NewKafkaConsumer creates a new kafka consumer
func NewKafkaConsumer(addr string, groupID string) (Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": addr,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	return &kafkaConsumerClient{svc: consumer}, nil
}

// Consume consumes messages from the kafka topic
func (c *kafkaConsumerClient) Consume(ctx context.Context, topic string, event chan<- *EventWrapper) {
	err := c.svc.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return
	}

	for {
		select {
		case <-ctx.Done():
			close(event)
			return
		default:
			msg, err := c.svc.ReadMessage(-1)
			if err != nil {
				log.WithError(err).Errorf("Consumer error: (%v)\n", msg)
				return
			}

			e := &EventWrapper{}
			err = json.Unmarshal(msg.Value, e)
			if err != nil {
				log.WithError(err).Errorf("failed to unmarshal message: (%v)\n", string(msg.Value))
				return
			}

			e.Topic = *msg.TopicPartition.Topic
			event <- e
		}
	}
}

// Close closes the kafka consumer
func (c *kafkaConsumerClient) Close() error {
	return c.svc.Close()
}
