package bullhorn

import (
	_ "embed"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

//go:embed test_data/valid_sub_response.json
var validSub string

//go:embed test_data/invalid_sub_response.json
var invalidSub string

var (
	okServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = fmt.Fprint(w, validSub)
	}))
	invalidResponseServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = fmt.Fprint(w, invalidSub)
	}))
	nonOkServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	}))
)

func TestNew(t *testing.T) {
	tests := []struct {
		name            string
		subscriptionUrl string
		entityUrl       string
	}{
		{
			name:            "can create new bullhorn instance",
			subscriptionUrl: "https://bullhorn.com/subscription",
			entityUrl:       "https://bullhorn.com/entity",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.subscriptionUrl, tt.entityUrl, &retryablehttp.Client{})
			assert.NotNil(t, got)
		})
	}
}

func TestClient_GetEvents(t *testing.T) {
	tests := []struct {
		name            string
		subscriptionUrl string
		entityUrl       string
		server          *httptest.Server
		oauthConfig     *oauth2.Config
		want            []Event
		wantErr         bool
	}{
		{
			name:            "ok",
			server:          okServer,
			subscriptionUrl: okServer.URL,
			entityUrl:       "https://bullhorn.com/entity",
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
			name:            "non OK response",
			server:          nonOkServer,
			subscriptionUrl: nonOkServer.URL,
			entityUrl:       "https://bullhorn.com/entity",
			want:            nil,
			wantErr:         true,
		},
		{
			name:            "invalid response",
			server:          invalidResponseServer,
			subscriptionUrl: nonOkServer.URL,
			entityUrl:       "https://bullhorn.com/entity",
			want:            nil,
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()

			c := New(tt.subscriptionUrl, tt.entityUrl, &retryablehttp.Client{
				CheckRetry: retryablehttp.DefaultRetryPolicy,
				Backoff:    retryablehttp.DefaultBackoff,
			})
			assert.NotNil(t, c)

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
