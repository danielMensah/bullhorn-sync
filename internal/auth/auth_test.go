package auth

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		password      string
		auth          *Mock
		oauthConfig   oauth2.Config
		expectMocks   func(t *testing.T, auth *Mock)
		expectedError string
	}{
		{
			name:     "successfully returns auth instance",
			username: "usr",
			password: "pass",
			auth:     &Mock{},
			oauthConfig: oauth2.Config{
				ClientID:     "clientID",
				ClientSecret: "secret",
				Endpoint: oauth2.Endpoint{
					AuthURL:   "authUrl",
					TokenURL:  "tokeUrl",
					AuthStyle: 0,
				},
				RedirectURL: "redirectUrl",
				Scopes:      nil,
			},
			expectMocks: func(t *testing.T, auth *Mock) {
				action := oauth2.SetAuthURLParam("action", "Login")
				username := oauth2.SetAuthURLParam("username", "usr")
				password := oauth2.SetAuthURLParam("password", "pass")

				ctx := context.TODO()
				token := &oauth2.Token{}
				l := []oauth2.AuthCodeOption{action, username, password}
				var d []oauth2.AuthCodeOption

				auth.On("AuthCodeURL", mock.Anything, l).Return("code")
				auth.On("Exchange", ctx, mock.Anything, d).Return(token, nil)
				auth.On("Client", ctx, token).Return(&http.Client{})
			},
		},
		{
			name:     "exchange error",
			username: "usr",
			password: "pass",
			auth:     &Mock{},
			oauthConfig: oauth2.Config{
				ClientID:     "clientID",
				ClientSecret: "secret",
				Endpoint: oauth2.Endpoint{
					AuthURL:   "authUrl",
					TokenURL:  "tokeUrl",
					AuthStyle: 0,
				},
				RedirectURL: "redirectUrl",
				Scopes:      nil,
			},
			expectMocks: func(t *testing.T, auth *Mock) {
				action := oauth2.SetAuthURLParam("action", "Login")
				username := oauth2.SetAuthURLParam("username", "usr")
				password := oauth2.SetAuthURLParam("password", "pass")

				ctx := context.TODO()
				l := []oauth2.AuthCodeOption{action, username, password}
				var d []oauth2.AuthCodeOption

				auth.On("AuthCodeURL", mock.Anything, l).Return("code")
				auth.On("Exchange", ctx, mock.Anything, d).Return(nil, errors.New("boom"))
			},
			expectedError: "exchanging code for token: boom",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectMocks != nil {
				tt.expectMocks(t, tt.auth)
			}

			ctx := context.TODO()
			client, err := New(ctx, tt.username, tt.password, tt.auth)

			if tt.expectedError == "" {
				assert.NoError(t, err)
				assert.NotNil(t, client)
			} else {
				assert.Contains(t, err.Error(), tt.expectedError)
			}

			if tt.expectMocks != nil {
				tt.auth.AssertExpectations(t)
			}
		})
	}
}
