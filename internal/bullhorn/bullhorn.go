package bullhorn

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/danielMensah/bullhorn-sync-poc/internal/auth"
	"github.com/danielMensah/bullhorn-sync-poc/internal/config"
	pb "github.com/danielMensah/bullhorn-sync-poc/internal/proto"
	"github.com/hashicorp/go-retryablehttp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Client is the Bullhorn client
type Client struct {
	subscriptionUrl string
	entityUrl       string
	httpClient      *retryablehttp.Client
}

type Bullhorn interface {
	GetEvents() ([]*pb.Event, error)
	FetchEntityChanges(event *pb.Event) (*Entity, error)
}

// New returns a new Bullhorn client.
func New(ctx context.Context, config *config.Config) (Bullhorn, error) {
	oauthConfig := &oauth2.Config{
		ClientID:     config.BullhornClientID,
		ClientSecret: config.BullhornSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   config.BullhornAuthUrl,
			TokenURL:  config.BullhornTokenUrl,
			AuthStyle: 0,
		},
		RedirectURL: config.BullhornRedirectUrl,
		Scopes:      nil,
	}

	client, err := auth.New(ctx, config.BullhornUsername, config.BullhornPassword, oauthConfig)
	if err != nil {
		return nil, err
	}

	return &Client{
		subscriptionUrl: config.BullhornSubscriptionUrl,
		entityUrl:       config.BullhornEntityUrl,
		httpClient:      client,
	}, nil
}

// GetEvents returns the events from the bullhorn subscription url
func (c *Client) GetEvents() ([]*pb.Event, error) {
	body, err := c.request("GET", c.subscriptionUrl, nil)
	if err != nil {
		return nil, err
	}

	var response RequestResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}

	events := make([]*pb.Event, 0)
	for _, event := range response.Events {
		ts := timestamppb.New(time.UnixMilli(event.EventTimestamp))

		events = append(events, &pb.Event{
			EntityId:          event.EntityId,
			EntityName:        event.EntityName,
			EntityEventType:   event.EntityEventType,
			UpdatedProperties: event.UpdatedProperties,
			EventTimestamp:    ts,
		})
	}

	return events, nil
}

// FetchEntityChanges fetches the changes for a given entity
func (c *Client) FetchEntityChanges(event *pb.Event) (*Entity, error) {
	switch event.EntityEventType {
	case pb.EventType_UPDATED:

		fields := strings.Join(event.UpdatedProperties, ",")
		url := fmt.Sprintf("%s/entity/%s/%d?fields=%s", c.entityUrl, event.EntityName, event.EntityId, fields)

		body, err := c.request("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed getting entity (%s) with id: %d : %w", event.EntityName, event.EntityId, err)
		}

		return &Entity{
			Id:        event.EntityId,
			Name:      event.EntityName,
			Changes:   body,
			Timestamp: event.EventTimestamp,
		}, nil
	case pb.EventType_INSERTED:
		fields := strings.Join(event.UpdatedProperties, ",")
		url := fmt.Sprintf("%s/entity/%s/%d?fields=%s", c.entityUrl, event.EntityName, event.EntityId, fields)

		body, err := c.request("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed getting entity (%s) with id: %d : %w", event.EntityName, event.EntityId, err)
		}

		return &Entity{
			Id:        event.EntityId,
			Name:      event.EntityName,
			Changes:   body,
			Timestamp: event.EventTimestamp,
		}, nil
	case pb.EventType_DELETED:
		// delete event entity
	default:
		log.Errorf("entity event type not supported: %s", event.EntityEventType)
	}

	return nil, nil
}
