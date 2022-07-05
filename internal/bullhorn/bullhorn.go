package bullhorn

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
)

// Client is the Bullhorn client
type Client struct {
	subscriptionUrl string
	entityUrl       string
	httpClient      *retryablehttp.Client
}

// Bullhorn is the interface that defines the methods that a bullhorn client must implement
type Bullhorn interface {
	GetEvents() ([]Event, error)
	FetchEntityChanges(event Event) (Entity, error)
	request(url string) ([]byte, error)
}

// New returns a new Bullhorn client.
func New(subscriptionUrl, entityUrl string, httpClient *retryablehttp.Client) Bullhorn {
	return &Client{
		subscriptionUrl: subscriptionUrl,
		entityUrl:       entityUrl,
		httpClient:      httpClient,
	}
}

// GetEvents returns the events from the bullhorn subscription url
func (c *Client) GetEvents() ([]Event, error) {
	body, err := c.request(c.subscriptionUrl)
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

		body, err := c.request(url)
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

		body, err := c.request(url)
		if err != nil {
			return Entity{}, fmt.Errorf("failed getting entity (%s) with id: %d : %w", event.EntityName, event.EntityId, err)
		}

		return Entity{
			Id:        event.EntityId,
			Name:      event.EntityName,
			EventType: event.EntityEventType,
			Changes:   body,
			Timestamp: event.EventTimestamp,
		}, nil
	case EventTypeDeleted:
		return Entity{
			Id:        event.EntityId,
			Name:      event.EntityName,
			EventType: event.EntityEventType,
			Timestamp: event.EventTimestamp,
		}, nil
	}

	return Entity{}, fmt.Errorf("entity event type not supported: %s", event.EntityEventType)
}

func (c *Client) request(url string) ([]byte, error) {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading request body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non ok response for request: %d status code %v: %s", resp.StatusCode, resp, string(response))
	}

	return response, nil
}
