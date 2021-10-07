package bullhorn

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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
		name        string
		config      Config
		oauthConfig *oauth2.Config
		wantErr     bool
	}{
		{
			name: "ok",
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
			wantErr: false,
		},
		{
			name: "handle error",
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
					TokenURL:  "someUrl",
					AuthStyle: 0,
				},
				RedirectURL: "redirect",
				Scopes:      nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			got, err := New(ctx, tt.config, tt.oauthConfig)

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
		config      Config
		oauthConfig *oauth2.Config
		want        []Event
		wantErr     bool
	}{
		{
			name: "ok",
			config: Config{
				Username:        "user",
				Password:        "pass",
				SubscriptionUrl: validSubServer.URL,
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
			config: Config{
				Username:        "user",
				Password:        "pass",
				SubscriptionUrl: invalidSubServer.URL,
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
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid response from sub call",
			config: Config{
				Username:        "user",
				Password:        "pass",
				SubscriptionUrl: invalidSubServer.URL,
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
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New(context.TODO(), tt.config, tt.oauthConfig)
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
