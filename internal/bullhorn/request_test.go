package bullhorn

import (
	"context"
	"io"
	"testing"

	"github.com/danielMensah/bullhorn-sync-poc/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestClient_request(t *testing.T) {
	tests := []struct {
		name             string
		config           *config.Config
		url              string
		body             io.Reader
		expectedResponse string
		wantErr          bool
	}{
		{
			name: "ok response",
			config: &config.Config{
				BullhornUsername:        "user",
				BullhornPassword:        "pass",
				BullhornSubscriptionUrl: "sub",
				BullhornEntityUrl:       "ent",
			},
			url:              validSubServer.URL,
			body:             nil,
			expectedResponse: validSub,
			wantErr:          false,
		},
		{
			name: "non ok response",
			config: &config.Config{
				BullhornUsername:        "user",
				BullhornPassword:        "pass",
				BullhornSubscriptionUrl: "sub",
				BullhornEntityUrl:       "ent",
			},
			url:              invalidSubServer.URL,
			body:             nil,
			expectedResponse: "",
			wantErr:          true,
		},
		{
			name: "request error",
			config: &config.Config{
				BullhornUsername:        "user",
				BullhornPassword:        "pass",
				BullhornSubscriptionUrl: "sub",
				BullhornEntityUrl:       "ent",
			},
			url:              "",
			body:             nil,
			expectedResponse: "",
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New(context.TODO(), tt.config)
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
