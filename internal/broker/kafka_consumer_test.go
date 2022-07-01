package broker

import (
	"context"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// KafkaPublisherMock is a mock of the queue methods for use in the services using the broker client
type KafkaReaderMock struct {
	mock.Mock
}

func (m *KafkaReaderMock) ReadMessage(ctx context.Context) (kafka.Message, error) {
	args := m.Called(ctx)
	message, ok := args.Get(0).(kafka.Message)
	if !ok {
		message = kafka.Message{}
	}

	return message, args.Error(1)
}

func (m *KafkaReaderMock) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestKafkaConsumerClient_Consume(t *testing.T) {
	tests := []struct {
		name          string
		reader        *KafkaReaderMock
		expectedEvent *EventWrapper
		expectMocks   func(t *testing.T, reader *KafkaReaderMock)
	}{
		{
			name:   "can consume event",
			reader: &KafkaReaderMock{},
			expectedEvent: &EventWrapper{
				Topic: "test",
				Data:  "some data",
			},
			expectMocks: func(t *testing.T, reader *KafkaReaderMock) {
				message := kafka.Message{
					Value: []byte(`{"topic":"test","data":"some data"}`),
				}

				reader.On("ReadMessage", mock.Anything).Return(message, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectMocks != nil {
				tt.expectMocks(t, tt.reader)
			}

			ctx, cancel := context.WithCancel(context.Background())

			s := &KafkaConsumerClient{
				reader: tt.reader,
			}

			events := make(chan *EventWrapper)

			go func() {
				for event := range events {
					assert.Equal(t, tt.expectedEvent, event)
					cancel()
				}
			}()

			s.Consume(ctx, events)

			if tt.expectMocks != nil {
				tt.reader.AssertExpectations(t)
			}
		})
	}
}
