package broker

import (
	"context"
	"encoding/json"
	"sync"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// KafkaPublisherMock is a mock of the queue methods for use in the services using the broker client
type KafkaPublisherMock struct {
	mock.Mock
}

func (m *KafkaPublisherMock) WriteMessages(msgs ...kafka.Message) (int, error) {
	args := m.Called(msgs)
	return 0, args.Error(1)
}

func (m *KafkaPublisherMock) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestKafkaPublisherClient_Publish(t *testing.T) {
	tests := []struct {
		name        string
		workers     int
		data        []*EventWrapper
		conn        *KafkaPublisherMock
		expectMocks func(t *testing.T, conn *KafkaPublisherMock)
	}{
		{
			name:    "can publish single event",
			conn:    &KafkaPublisherMock{},
			workers: 1,
			data: []*EventWrapper{
				{
					Topic: "test",
					Data:  "some data 1",
				},
			},
			expectMocks: func(t *testing.T, conn *KafkaPublisherMock) {
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
			conn:    &KafkaPublisherMock{},
			workers: 10,
			data: []*EventWrapper{
				{
					Topic: "test",
					Data:  "some data 1",
				},
			},
			expectMocks: func(t *testing.T, conn *KafkaPublisherMock) {
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
			conn:    &KafkaPublisherMock{},
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
			expectMocks: func(t *testing.T, conn *KafkaPublisherMock) {
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
			conn:    &KafkaPublisherMock{},
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
			expectMocks: func(t *testing.T, conn *KafkaPublisherMock) {
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

			c := &KafkaPublisherClient{
				Conn: tt.conn,
			}

			events := make(chan *EventWrapper)
			wg := &sync.WaitGroup{}

			for i := 0; i < tt.workers; i++ {
				wg.Add(1)
				go c.Publish(context.Background(), events, wg)
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
