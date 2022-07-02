package broker

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// KafkaConsumerMock is a mock of the queue methods for use in the services using the broker client
type KafkaConsumerMock struct {
	mock.Mock
}

func (m *KafkaConsumerMock) ReadMessage(timeout time.Duration) (*kafka.Message, error) {
	args := m.Called(timeout)
	message, ok := args.Get(0).(*kafka.Message)
	if !ok {
		message = &kafka.Message{}
	}

	return message, args.Error(1)
}

func (m *KafkaConsumerMock) SubscribeTopics(topics []string, rebalanceCb kafka.RebalanceCb) (err error) {
	args := m.Called(topics, rebalanceCb)
	return args.Error(0)
}

func (m *KafkaConsumerMock) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestKafkaConsumerClient_Consume(t *testing.T) {
	tests := []struct {
		name          string
		topic         string
		consumerMock  *KafkaConsumerMock
		expectedEvent *EventWrapper
		expectMocks   func(t *testing.T, consumerMock *KafkaConsumerMock)
	}{
		{
			name:         "can consume event",
			topic:        "test",
			consumerMock: &KafkaConsumerMock{},
			expectedEvent: &EventWrapper{
				Topic: "test",
				Event: "some data",
			},
			expectMocks: func(t *testing.T, consumerMock *KafkaConsumerMock) {
				topic := "test"
				msg := &kafka.Message{
					TopicPartition: kafka.TopicPartition{
						Topic:     &topic,
						Partition: 0,
					},
					Value: []byte(`{"event":"some data"}`),
				}

				consumerMock.On("SubscribeTopics", []string{topic}, mock.AnythingOfType("kafka.RebalanceCb")).Return(nil)
				consumerMock.On("ReadMessage", mock.AnythingOfType("time.Duration")).Return(msg, nil)
				consumerMock.On("Close").Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectMocks != nil {
				tt.expectMocks(t, tt.consumerMock)
			}

			ctx, cancel := context.WithCancel(context.Background())

			s := &KafkaConsumerClient{
				svc: tt.consumerMock,
			}

			events := make(chan *EventWrapper)
			wg := &sync.WaitGroup{}

			wg.Add(1)
			go func() {
				for event := range events {
					assert.Equal(t, tt.expectedEvent, event)
					cancel()
				}

				wg.Done()
			}()

			s.Consume(ctx, tt.topic, events)
			wg.Wait()

			if tt.expectMocks != nil {
				tt.consumerMock.AssertExpectations(t)
			}
		})
	}
}
