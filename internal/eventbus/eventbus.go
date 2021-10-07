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

const (
	ProcessedSource = "bullhorn-sync.processed"
	IngestionSource = "bullhorn-sync.ingestion"
)

// Client for queue
type Client struct {
	eventbus eventbridgeiface.EventBridgeAPI
}

type loggerResults struct {
	successful int
	failed     int
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

// SendIngestedEvents puts events to eventBridge
func (c Client) SendIngestedEvents(ctx context.Context, source string, events []bullhorn.Event) error {
	var entries []*eventbridge.PutEventsRequestEntry

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
			EventBusName: aws.String("bullhorn-sync"),
			Source:       aws.String(source),
			DetailType:   aws.String(event.EntityName),
			Detail:       aws.String(string(payload)),
			Time:         &t,
		})
	}

	p, err := c.eventbus.PutEventsWithContext(ctx, &eventbridge.PutEventsInput{Entries: entries})
	if err != nil {
		return fmt.Errorf("putting events into eventBridge: %w", err)
	}

	res := loggerResults{}
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

// SendProcessedEvent puts events to eventBridge
func (c Client) SendProcessedEvent(ctx context.Context, source string, entity bullhorn.Entity) error {
	log.WithFields(log.Fields{
		"entityId":   entity.Id,
		"entityName": entity.Name,
	}).Info("processing")

	payload, err := json.Marshal(entity)
	if err != nil {
		return fmt.Errorf("marshalling event: %w", err)
	}

	t := time.Unix(entity.Timestamp, 0).UTC()
	entry := eventbridge.PutEventsRequestEntry{
		EventBusName: aws.String("bullhorn-sync"),
		Source:       aws.String(source),
		DetailType:   aws.String(entity.Name),
		Detail:       aws.String(string(payload)),
		Time:         &t,
	}

	p, err := c.eventbus.PutEventsWithContext(ctx, &eventbridge.PutEventsInput{
		Entries: []*eventbridge.PutEventsRequestEntry{&entry},
	})
	if err != nil {
		return fmt.Errorf("putting events into eventBridge: %w", err)
	}

	for _, outputEntry := range p.Entries {
		if outputEntry.EventId != nil {
			log.Info("successfully sent event to event bus")
		} else {
			log.Warn("failed to send event to event bus")
		}
	}

	return nil
}
