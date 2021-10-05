package auth

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

// Mock is a mock of the queue methods for use in the services using the queue client
type Mock struct {
	mock.Mock
}

// AuthCodeURL is a mock of this method used for testing
func (m *Mock) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	args := m.Called(state, opts)
	code, ok := args.Get(0).(string)
	if !ok {
		code = ""
	}
	return code
}

// Exchange is a mock of this method used for testing
func (m *Mock) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	args := m.Called(ctx, code, opts)
	token, ok := args.Get(0).(*oauth2.Token)
	if !ok {
		token = nil
	}
	return token, args.Error(1)
}

// Client is a mock of this method used for testing
func (m *Mock) Client(ctx context.Context, t *oauth2.Token) *http.Client {
	args := m.Called(ctx, t)
	client, ok := args.Get(0).(*http.Client)
	if !ok {
		client = nil
	}
	return client
}
