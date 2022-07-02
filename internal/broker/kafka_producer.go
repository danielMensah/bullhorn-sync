package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
)

type KafkaProducerClient struct {
	svc KafkaProducerService
}

// KafkaProducerService is for mainly testing purposes
type KafkaProducerService interface {
	ProduceChannel() chan *kafka.Message
	Events() chan kafka.Event
	Flush(timeoutMs int) int
	Len() int
	Close()
}

func NewKafkaProducer(addr string) (Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": addr})
	if err != nil {
		return nil, err
	}

	p.Events()
	return &KafkaProducerClient{svc: p}, nil
}

func (p *KafkaProducerClient) Produce(ctx context.Context, topic string, events <-chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			for p.svc.Flush(100) > 0 {
				log.Warnf("%d messages still pending", p.svc.Len())
			}

			p.svc.Close()
			return
		case event, ok := <-events:
			if !ok {
				return
			}

			marshalledData, err := json.Marshal(event)
			if err != nil {
				log.WithError(err).Error("marshalling event")
				return
			}

			message := &kafka.Message{
				Value:          marshalledData,
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			}

			p.svc.ProduceChannel() <- message
		}
	}
}

func (p *KafkaProducerClient) MonitorEvents(wg *sync.WaitGroup) {
	defer wg.Done()

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
