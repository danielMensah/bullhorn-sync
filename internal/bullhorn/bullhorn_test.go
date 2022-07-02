package bullhorn

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danielMensah/bullhorn-sync-poc/internal/config"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

const (
	validSub   = `{"requestId":1,"events":[{"eventId":"abc","eventTimestamp":1633475158,"entityName":"candidate","entityId":1,"entityEventType":"UPDATE","updatedProperties":["name","dob"]}]}`
	invalidSub = `"requestId":1,"events":[{"eventId":"abc","eventTimestamp":1633475158,"entityName":"candidate","entityId":1,"entityEventType":"UPDATE","updatedProperties":["name","dob"]}]`
)

var (
	tokenServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = fmt.Fprint(w, "access_token=eeybb")
	}))
	validSubServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = fmt.Fprint(w, validSub)
	}))
	invalidSubServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		_, _ = fmt.Fprint(w, invalidSub)
	}))
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  *config.Config
		wantErr bool
	}{
		{
			name: "ok",
			config: &config.Config{
				BullhornUsername:        "user",
				BullhornPassword:        "pass",
				BullhornSubscriptionUrl: "sub",
				BullhornEntityUrl:       "ent",
			},
			wantErr: false,
		},
		{
			name: "handle error",
			config: &config.Config{
				BullhornUsername:        "user",
				BullhornPassword:        "pass",
				BullhornSubscriptionUrl: "sub",
				BullhornEntityUrl:       "ent",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			got, err := New(ctx, tt.config)

			if !tt.wantErr {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestClient_GetEvents(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.Config
		oauthConfig *oauth2.Config
		want        []Event
		wantErr     bool
	}{
		{
			name: "ok",
			config: &config.Config{
				BullhornUsername:        "user",
				BullhornPassword:        "pass",
				BullhornSubscriptionUrl: validSubServer.URL,
			},
			want: []Event{
				{
					EventId:           "abc",
					EventTimestamp:    int64(1633475158),
					EntityName:        "candidate",
					EntityId:          1,
					EntityEventType:   "UPDATE",
					UpdatedProperties: []string{"name", "dob"},
				},
			},
			wantErr: false,
		},
		{
			name: "non OK response",
			config: &config.Config{
				BullhornUsername:        "user",
				BullhornPassword:        "pass",
				BullhornSubscriptionUrl: invalidSubServer.URL,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid response from sub call",
			config: &config.Config{
				BullhornUsername:        "user",
				BullhornPassword:        "pass",
				BullhornSubscriptionUrl: invalidSubServer.URL,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New(context.Background(), tt.config)
			assert.NoError(t, err)

			got, err := c.GetEvents()

			if !tt.wantErr {
				assert.NoError(t, err)
				assert.ElementsMatch(t, tt.want, got)
			} else {
				assert.Error(t, err)
			}

		})
	}
}
