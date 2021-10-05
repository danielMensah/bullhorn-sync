package bullhorn

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/danielMensah/bullhorn-sync-poc/internal/auth"
	"github.com/hashicorp/go-retryablehttp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
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
	retryClient     *retryablehttp.Client
}

//oauthClient := oauth2.Config{
//ClientID:     clientID,
//ClientSecret: clientSecret,
//Endpoint: oauth2.Endpoint{
//AuthURL:   "https://auth-emea.bullhornstaffing.com/oauth/authorize",
//TokenURL:  "https://rest-emea.bullhornstaffing.com/rest-services/login?version=2.0",
//AuthStyle: 0,
//},
//RedirectURL: "http://www.bullhorn.com",
//Scopes:      nil,
//}

func New(ctx context.Context, config Config, oauthConfig *oauth2.Config) *Client {
	client, err := auth.New(ctx, config.Username, config.Password, oauthConfig)
	if err != nil {

	}

	return &Client{
		subscriptionUrl: config.SubscriptionUrl,
		entityUrl:       config.EntityUrl,
		retryClient:     client,
	}
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

func (c *Client) FetchChanges(event Event, record chan<- Record) {
	switch event.EntityEventType {
	case updated:
		fields := strings.Join(event.UpdatedProperties, ",")
		url := fmt.Sprintf("%s/entity/%s/%d?fields=%s", c.entityUrl, event.EntityName, event.EntityId, fields)

		body, err := c.request("GET", url, nil)
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
