package bullhorn

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/danielMensah/bullhorn-sync-poc/internal/auth"
	"github.com/hashicorp/go-retryablehttp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

var (
	updated  = "UPDATED"
	inserted = "INSERTED"
	deleted  = "DELETED"
)

type Config struct {
	Username        string
	Password        string
	SubscriptionUrl string
	EntityUrl       string
}

type Client struct {
	subscriptionUrl string
	entityUrl       string
	httpClient      *retryablehttp.Client
}

func New(ctx context.Context, config Config, oauthConfig *oauth2.Config) (*Client, error) {
	client, err := auth.New(ctx, config.Username, config.Password, oauthConfig)
	if err != nil {
		return nil, err
	}

	return &Client{
		subscriptionUrl: config.SubscriptionUrl,
		entityUrl:       config.EntityUrl,
		httpClient:      client,
	}, nil
}

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

func (c *Client) FetchEntityChanges(event Event) (Entity, error) {
	switch event.EntityEventType {
	case updated:
		fields := strings.Join(event.UpdatedProperties, ",")
		url := fmt.Sprintf("%s/entity/%s/%d?fields=%s", c.entityUrl, event.EntityName, event.EntityId, fields)

		body, err := c.request("GET", url, nil)
		if err != nil {
			return Entity{}, fmt.Errorf("failed getting entity (%s) with id: %d : %w", event.EntityName, event.EntityId, err)
		}

		return Entity{
			Id:        event.EntityId,
			Name:      event.EntityName,
			Changes:   string(body),
			Timestamp: event.EventTimestamp,
		}, nil
	case inserted:
		fields := strings.Join(event.UpdatedProperties, ",")
		url := fmt.Sprintf("%s/entity/%s/%d?fields=%s", c.entityUrl, event.EntityName, event.EntityId, fields)

		body, err := c.request("GET", url, nil)
		if err != nil {
			return Entity{}, fmt.Errorf("failed getting entity (%s) with id: %d : %w", event.EntityName, event.EntityId, err)
		}

		return Entity{
			Id:        event.EntityId,
			Name:      event.EntityName,
			Changes:   string(body),
			Timestamp: event.EventTimestamp,
		}, nil
	case deleted:
		// TODO
	default:
		log.Errorf("entity event type not supported: %s", event.EntityEventType)
	}

	return Entity{}, nil
}
