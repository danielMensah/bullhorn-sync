package bullhorn

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/danielMensah/bullhorn-sync-poc/internal/auth"
	"github.com/danielMensah/bullhorn-sync-poc/internal/config"
	"github.com/hashicorp/go-retryablehttp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// Client is the Bullhorn client
type Client struct {
	subscriptionUrl string
	entityUrl       string
	httpClient      *retryablehttp.Client
}

type Bullhorn interface {
	GetEvents() ([]Event, error)
	FetchEntityChanges(event Event) (Entity, error)
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
func (c *Client) GetEvents() ([]Event, error) {
	body, err := c.request("GET", c.subscriptionUrl, nil)
	if err != nil {
		return nil, err
	}

	var response RequestResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}

	return response.Events, nil
}

// FetchEntityChanges fetches the changes for a given entity
func (c *Client) FetchEntityChanges(event Event) (Entity, error) {
	switch event.EntityEventType {
	case EventTypeInserted:
		fields := "*"
		url := fmt.Sprintf("%s/entity/%s/%d?fields=%s", c.entityUrl, event.EntityName, event.EntityId, fields)

		body, err := c.request("GET", url, nil)
		if err != nil {
			return Entity{}, fmt.Errorf("failed getting entity (%s) with id: %d : %w", event.EntityName, event.EntityId, err)
		}

		return Entity{
			Id:        event.EntityId,
			Name:      event.EntityName,
			Changes:   body,
			Timestamp: event.EventTimestamp,
		}, nil
	case EventTypeUpdated:

		fields := strings.Join(event.UpdatedProperties, ",")
		url := fmt.Sprintf("%s/entity/%s/%d?fields=%s", c.entityUrl, event.EntityName, event.EntityId, fields)

		body, err := c.request("GET", url, nil)
		if err != nil {
			return Entity{}, fmt.Errorf("failed getting entity (%s) with id: %d : %w", event.EntityName, event.EntityId, err)
		}

		return Entity{
			Id:        event.EntityId,
			Name:      event.EntityName,
			EventType: string(event.EntityEventType),
			Changes:   body,
			Timestamp: event.EventTimestamp,
		}, nil
	case EventTypeDeleted:
		// delete event entity
	default:
		log.Errorf("entity event type not supported: %s", event.EntityEventType)
	}

	return Entity{}, nil
}
