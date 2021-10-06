package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/aws/aws-sdk-go/service/eventbridge/eventbridgeiface"
	"github.com/danielMensah/bullhorn-sync-poc/internal/bullhorn"
	log "github.com/sirupsen/logrus"
)

// Client for queue
type Client struct {
	eventbus eventbridgeiface.EventBridgeAPI
}

// Bus is an interface to expose eventbus methods
type Bus interface {
	Put(ctx context.Context, events []bullhorn.Event) error
}

// New will create a EventBridge client with a session
func New(region string) *Client {
	bus := eventbridge.New(session.Must(session.NewSession()), &aws.Config{
		Region: &region,
	})

	return &Client{eventbus: bus}
}

// Put puts events to eventBridge
func (c Client) Put(ctx context.Context, events []bullhorn.Event) error {
	var entries []*eventbridge.PutEventsRequestEntry

	res := &struct {
		successful int
		failed     int
	}{}

	for _, event := range events {
		log.WithFields(log.Fields{
			"entityId":   event.EntityId,
			"entityName": event.EntityName,
			"actionType": event.EntityEventType,
		}).Info("processing")

		payload, err := json.Marshal(event)
		if err != nil {
			log.WithError(err).Error("marshalling event")
			continue
		}

		t := time.Unix(event.EventTimestamp, 0).UTC()
		entries = append(entries, &eventbridge.PutEventsRequestEntry{
			Detail:       aws.String(string(payload)),
			DetailType:   aws.String(event.EntityEventType),
			EventBusName: aws.String("bullhorn-sync"),
			Source:       aws.String("bullhorn-sync.ingestion"),
			Time:         &t,
		})
	}

	p, err := c.eventbus.PutEventsWithContext(ctx, &eventbridge.PutEventsInput{Entries: entries})
	if err != nil {
		return fmt.Errorf("putting events into eventBridge: %w", err)
	}

	for _, outputEntry := range p.Entries {
		if outputEntry.EventId != nil {
			res.successful += 1
		} else {
			res.failed += 1
		}
	}

	log.WithFields(log.Fields{
		"successful": res.successful,
		"failed":     res.failed,
	}).Info("complete")
	return nil
}
