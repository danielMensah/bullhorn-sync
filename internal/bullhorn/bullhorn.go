package bullhorn

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/danielMensah/bullhorn-sync-poc/internal/auth"
	"github.com/hashicorp/go-retryablehttp"
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
