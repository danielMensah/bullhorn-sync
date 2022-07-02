package broker

import (
	"context"
	"encoding/json"
	"sync"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// KafkaProducerMock is a mock of the queue methods for use in the services using the broker client
type KafkaProducerMock struct {
	mock.Mock
}

func (m *KafkaProducerMock) ProduceChannel() chan *kafka.Message {
	args := m.Called()
	msg := args.Get(0).(chan *kafka.Message)
	if msg != nil {
		msg = make(chan *kafka.Message)
	}

	return msg
}

func (m *KafkaProducerMock) Events() chan kafka.Event {
	args := m.Called()
	event := args.Get(0).(chan kafka.Event)
	if event != nil {
		event = make(chan kafka.Event)
	}

	return event
}

func (m *KafkaProducerMock) Close() {
	_ = m.Called()
}

func TestKafkaPublisherClient_Publish(t *testing.T) {
	tests := []struct {
		name        string
		workers     int
		data        []*EventWrapper
		conn        *KafkaProducerMock
		expectMocks func(t *testing.T, conn *KafkaProducerMock)
	}{
		{
			name:    "can publish single event",
			conn:    &KafkaProducerMock{},
			workers: 1,
			data: []*EventWrapper{
				{
					Topic: "test",
					Data:  "some data 1",
				},
			},
			expectMocks: func(t *testing.T, conn *KafkaProducerMock) {
				d := &EventWrapper{Data: "some data 1"}

				v, err := json.Marshal(d.Data)
				assert.Nil(t, err)

				message := kafka.Message{
					Topic: "test",
					Value: v,
				}
				conn.On("WriteMessages", []kafka.Message{message}).Return(0, nil)
			},
		},
		{
			name:    "can publish single event with multiple goroutines",
			conn:    &KafkaProducerMock{},
			workers: 10,
			data: []*EventWrapper{
				{
					Topic: "test",
					Data:  "some data 1",
				},
			},
			expectMocks: func(t *testing.T, conn *KafkaProducerMock) {
				d := &EventWrapper{Data: "some data 1"}

				v, err := json.Marshal(d.Data)
				assert.Nil(t, err)

				message := kafka.Message{
					Topic: "test",
					Value: v,
				}
				conn.On("WriteMessages", []kafka.Message{message}).Return(0, nil)
			},
		},
		{
			name:    "can publish multiple events",
			conn:    &KafkaProducerMock{},
			workers: 1,
			data: []*EventWrapper{
				{
					Topic: "test",
					Data:  "some data 1",
				},
				{
					Topic: "test",
					Data:  "some data 2",
				},
				{
					Topic: "test",
					Data:  "some data 3",
				},
			},
			expectMocks: func(t *testing.T, conn *KafkaProducerMock) {
				d1 := &EventWrapper{Data: "some data 1"}
				d2 := &EventWrapper{Data: "some data 2"}
				d3 := &EventWrapper{Data: "some data 3"}

				v1, err := json.Marshal(d1.Data)
				assert.Nil(t, err)

				v2, err := json.Marshal(d2.Data)
				assert.Nil(t, err)

				v3, err := json.Marshal(d3.Data)
				assert.Nil(t, err)

				message1 := kafka.Message{
					Topic: "test",
					Value: v1,
				}

				message2 := kafka.Message{
					Topic: "test",
					Value: v2,
				}

				message3 := kafka.Message{
					Topic: "test",
					Value: v3,
				}

				conn.On("WriteMessages", []kafka.Message{message1}).Return(0, nil)
				conn.On("WriteMessages", []kafka.Message{message2}).Return(0, nil)
				conn.On("WriteMessages", []kafka.Message{message3}).Return(0, nil)
			},
		},
		{
			name:    "can publish multiple events with multiple goroutines",
			conn:    &KafkaProducerMock{},
			workers: 10,
			data: []*EventWrapper{
				{
					Topic: "test",
					Data:  "some data 1",
				},
				{
					Topic: "test",
					Data:  "some data 2",
				},
				{
					Topic: "test",
					Data:  "some data 3",
				},
			},
			expectMocks: func(t *testing.T, conn *KafkaProducerMock) {
				d1 := &EventWrapper{Data: "some data 1"}
				d2 := &EventWrapper{Data: "some data 2"}
				d3 := &EventWrapper{Data: "some data 3"}

				v1, err := json.Marshal(d1.Data)
				assert.Nil(t, err)

				v2, err := json.Marshal(d2.Data)
				assert.Nil(t, err)

				v3, err := json.Marshal(d3.Data)
				assert.Nil(t, err)

				message1 := kafka.Message{
					Topic: "test",
					Value: v1,
				}

				message2 := kafka.Message{
					Topic: "test",
					Value: v2,
				}

				message3 := kafka.Message{
					Topic: "test",
					Value: v3,
				}

				conn.On("WriteMessages", []kafka.Message{message1}).Return(0, nil)
				conn.On("WriteMessages", []kafka.Message{message2}).Return(0, nil)
				conn.On("WriteMessages", []kafka.Message{message3}).Return(0, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectMocks != nil {
				tt.expectMocks(t, tt.conn)
			}

			c := &KafkaProducerClient{
				svc: tt.conn,
			}

			events := make(chan *EventWrapper)
			wg := &sync.WaitGroup{}

			for i := 0; i < tt.workers; i++ {
				wg.Add(1)
				go c.Produce(context.Background(), events, wg)
			}

			for _, d := range tt.data {
				events <- d
			}

			close(events)
			wg.Wait()

			if tt.expectMocks != nil {
				tt.conn.AssertExpectations(t)
			}
		})
	}
}
