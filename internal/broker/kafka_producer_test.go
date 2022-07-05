package broker

import (
	"context"
	"sync"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/mock"
)

// KafkaProducerMock is a mock of the broker methods for use in the services using the broker client
type KafkaProducerMock struct {
	mock.Mock
}

func (m *KafkaProducerMock) ProduceChannel() chan *kafka.Message {
	args := m.Called()
	msg := args.Get(0).(chan *kafka.Message)

	go func() {
		<-msg
	}()

	return msg
}

func (m *KafkaProducerMock) Events() chan kafka.Event {
	args := m.Called()
	event := args.Get(0).(chan kafka.Event)

	return event
}

func (m *KafkaProducerMock) Flush(timeoutMs int) int {
	args := m.Called(timeoutMs)
	n := args.Get(0).(int)

	return n
}
func (m *KafkaProducerMock) Len() int {
	args := m.Called()
	n := args.Get(0).(int)

	return n
}

func (m *KafkaProducerMock) Close() {
	_ = m.Called()
}

func TestKafkaPublisherClient_Publish(t *testing.T) {
	tests := []struct {
		name         string
		topic        string
		event        []*EventWrapper
		workers      int
		producerMock *KafkaProducerMock
		expectMocks  func(t *testing.T, conn *KafkaProducerMock)
	}{
		{
			name:  "can produce event",
			topic: "test",
			event: []*EventWrapper{
				{
					Topic: "test",
					Event: "some data 1",
				},
				{
					Topic: "test",
					Event: "some data 2",
				},
			},
			workers:      1,
			producerMock: &KafkaProducerMock{},
			expectMocks: func(t *testing.T, producerMock *KafkaProducerMock) {
				producerMock.On("ProduceChannel").Return(make(chan *kafka.Message))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectMocks != nil {
				tt.expectMocks(t, tt.producerMock)
			}

			c := &kafkaProducerClient{
				svc: tt.producerMock,
			}

			events := make(chan *EventWrapper)
			wg := &sync.WaitGroup{}

			ctx, cancel := context.WithCancel(context.Background())
			for i := 0; i < tt.workers; i++ {
				wg.Add(1)
				go c.Produce(ctx, events, wg)
			}

			for _, event := range tt.event {
				events <- event
			}

			cancel()
			close(events)
			wg.Wait()

			if tt.expectMocks != nil {
				tt.producerMock.AssertExpectations(t)
			}
		})
	}
}
