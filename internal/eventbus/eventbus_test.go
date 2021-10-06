package eventbus

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/danielMensah/bullhorn-sync-poc/internal/bullhorn"
	"github.com/stretchr/testify/assert"
)

var (
	tt        = time.Date(2021, time.June, 21, 10, 53, 8, 0, time.UTC)
	putevents = &eventbridge.PutEventsInput{
		Entries: []*eventbridge.PutEventsRequestEntry{
			{
				Detail:       aws.String(`{"eventId":"371cdde3-4a24-4174-ba1c-c5cf6a4b672f","eventTimestamp":1624272788,"entityName":"Candidate","entityId":734,"entityEventType":"UPDATED","updatedProperties":["name","dob"]}`),
				DetailType:   aws.String("UPDATED"),
				EventBusName: aws.String("bullhorn-sync"),
				Source:       aws.String("bullhorn-sync.ingestion"),
				Time:         &tt,
			},
		},
	}
)

func TestNew(t *testing.T) {
	client := New("region")
	assert.NotNil(t, client)
}

func TestClient_Put(t *testing.T) {
	tests := []struct {
		name        string
		bus         *Mock
		events      []bullhorn.Event
		expectMocks func(t *testing.T, bus *Mock)
		expectedErr string
	}{
		{
			name: "ok",
			bus:  &Mock{},
			events: []bullhorn.Event{
				{
					EventId:           "371cdde3-4a24-4174-ba1c-c5cf6a4b672f",
					EventTimestamp:    tt.Unix(),
					EntityName:        "Candidate",
					EntityId:          734,
					EntityEventType:   "UPDATED",
					UpdatedProperties: []string{"name", "dob"},
				},
			},
			expectMocks: func(t *testing.T, bus *Mock) {
				failedCount := int64(0)
				eventId := "324dabee-331d-4251-9880-4b2eaa568968"
				output := &eventbridge.PutEventsOutput{
					Entries: []*eventbridge.PutEventsResultEntry{
						{
							EventId: &eventId,
						},
					},
					FailedEntryCount: &failedCount,
				}

				bus.On("PutEventsWithContext", context.TODO(), putevents).Return(output, nil)
			},
			expectedErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectMocks != nil {
				tt.expectMocks(t, tt.bus)
			}

			c := Client{eventbus: tt.bus}
			err := c.Put(context.TODO(), tt.events)

			if tt.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Contains(t, err.Error(), tt.expectedErr)
			}
		})
	}
}
