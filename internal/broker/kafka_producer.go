package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
)

type kafkaProducerClient struct {
	svc KafkaProducerService
}

// KafkaProducerService is an interface to mock the kafka producer for testing
type KafkaProducerService interface {
	// ProduceChannel returns the channel to send messages to the kafka topic
	ProduceChannel() chan *kafka.Message
	// Events returns the channel to listen for kafka events
	Events() chan kafka.Event
	// Flush flushes and wait for all messages to be delivered
	Flush(timeoutMs int) int
	// Len returns the number of messages pending in the producer
	Len() int
	// Close closes the producer
	Close()
}

// NewKafkaProducer creates a new KafkaProducerClient
func NewKafkaProducer(addr string) (Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": addr})
	if err != nil {
		return nil, err
	}

	return &kafkaProducerClient{svc: p}, nil
}

// Produce sends the events to the kafka topic
func (p *kafkaProducerClient) Produce(ctx context.Context, events <-chan *EventWrapper, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case wrappedEvent, ok := <-events:
			if !ok {
				return
			}

			marshalledEvent, err := json.Marshal(wrappedEvent.Event)
			if err != nil {
				log.WithError(err).Error("marshalling event")
				return
			}

			message := &kafka.Message{
				Value:          marshalledEvent,
				TopicPartition: kafka.TopicPartition{Topic: &wrappedEvent.Topic, Partition: kafka.PartitionAny},
			}

			p.svc.ProduceChannel() <- message
		}
	}
}

// MonitorEvents listens for kafka events and logs them
func (p *kafkaProducerClient) MonitorEvents() {
	for e := range p.svc.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			m := ev
			if m.TopicPartition.Error != nil {
				fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
			} else {
				fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
					*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
			}
			return

		default:
			fmt.Printf("Ignored event: %s\n", ev)
		}
	}
}

// Close closes the kafka producer
func (p *kafkaProducerClient) Close() {
	for p.svc.Flush(100) > 0 {
		log.Warnf("%d messages still pending", p.svc.Len())
	}

	p.svc.Close()
}
