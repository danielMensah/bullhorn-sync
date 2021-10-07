package bullhorn

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestClient_request(t *testing.T) {
	tests := []struct {
		name             string
		config           Config
		oauthConfig      *oauth2.Config
		url              string
		body             io.Reader
		expectedResponse string
		wantErr          bool
	}{
		{
			name: "ok response",
			config: Config{
				Username:        "user",
				Password:        "pass",
				SubscriptionUrl: "sub",
				EntityUrl:       "ent",
			},
			oauthConfig: &oauth2.Config{
				ClientID:     "clientId",
				ClientSecret: "secret",
				Endpoint: oauth2.Endpoint{
					AuthURL:   "auth",
					TokenURL:  tokenServer.URL,
					AuthStyle: 0,
				},
				RedirectURL: "redirect",
				Scopes:      nil,
			},
			url:              validSubServer.URL,
			body:             nil,
			expectedResponse: validSub,
			wantErr:          false,
		},
		{
			name: "non ok response",
			config: Config{
				Username:        "user",
				Password:        "pass",
				SubscriptionUrl: "sub",
				EntityUrl:       "ent",
			},
			oauthConfig: &oauth2.Config{
				ClientID:     "clientId",
				ClientSecret: "secret",
				Endpoint: oauth2.Endpoint{
					AuthURL:   "auth",
					TokenURL:  tokenServer.URL,
					AuthStyle: 0,
				},
				RedirectURL: "redirect",
				Scopes:      nil,
			},
			url:              invalidSubServer.URL,
			body:             nil,
			expectedResponse: "",
			wantErr:          true,
		},
		{
			name: "request error",
			config: Config{
				Username:        "user",
				Password:        "pass",
				SubscriptionUrl: "sub",
				EntityUrl:       "ent",
			},
			oauthConfig: &oauth2.Config{
				ClientID:     "clientId",
				ClientSecret: "secret",
				Endpoint: oauth2.Endpoint{
					AuthURL:   "auth",
					TokenURL:  tokenServer.URL,
					AuthStyle: 0,
				},
				RedirectURL: "redirect",
				Scopes:      nil,
			},
			url:              "",
			body:             nil,
			expectedResponse: "",
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New(context.TODO(), tt.config, tt.oauthConfig)
			assert.NoError(t, err)

			response, err := c.request("GET", tt.url, tt.body)

			if !tt.wantErr {
				assert.Equal(t, tt.expectedResponse, string(response))
			} else {
				assert.Error(t, err)
			}

		})
	}
}
