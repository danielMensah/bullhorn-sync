package bullhorn

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2/clientcredentials"
)

type Client struct {
	subUrl string
	base   string
	http   *http.Client
}

func New(subUrl, entityBaseUrl string) *Client {
	retryClient := retryablehttp.NewClient()
	retryClient.Logger = true

	oauthClient := clientcredentials.Config{
		ClientID:     cfg.WebhookIntegrationClientID,
		ClientSecret: cfg.WebhookIntegrationClientSecret,
		TokenURL:     fmt.Sprintf(`https://%s/oauth/token`, cfg.CustomDomain),
		EndpointParams: map[string][]string{
			"audience": {cfg.DefaultAudience},
		},
	}

	return &Client{
		subUrl: subUrl,
		base:   entityBaseUrl,
		http:   retryClient.HTTPClient,
	}
}

func (c *Client) GetEvents() ([]Event, error) {
	headers := map[string]string{
		"Authorization": "",
	}
	body, err := c.request("GET", c.subUrl, nil, headers)
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

func (c *Client) FetchChanges(event Event, record chan<- Record) {
	switch event.EntityEventType {
	case updated:
		fields := strings.Join(event.UpdatedProperties, ",")
		url := fmt.Sprintf("%s/entity/%s/%d?fields=%s", c.base, event.EntityName, event.EntityId, fields)

		headers := map[string]string{
			"Authorization": "",
		}

		body, err := c.request("GET", url, nil, headers)
		if err != nil {
			log.WithError(err).Error("getting updated fields")
			return
		}

		record <- Record{
			EntityId:        event.EntityId,
			EntityName:      event.EntityName,
			EntityEventType: event.EntityEventType,
			EventTimestamp:  time.Unix(event.EventTimestamp, 0),
			Changes:         body,
		}
	case inserted:
		// TODO
	case deleted:
		// TODO
	default:
		log.Errorf("entity event type not supported: %s", event.EntityEventType)
	}
}
