package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"golang.org/x/oauth2"
)

var (
	retryWaitMin = 1 * time.Second
	retryWaitMax = 30 * time.Second
	retryMax     = 4
	logger       = log.New(os.Stderr, "", log.LstdFlags)
)

// Oauth2 exposes oath2 methods
type Oauth2 interface {
	// AuthCodeURL returns a URL to OAuth 2.0 provider's consent page
	// that asks for permissions for the required scopes explicitly.
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	// Exchange converts an authorization code into a token.
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	// Client returns an HTTP client using the provided token.
	Client(ctx context.Context, t *oauth2.Token) *http.Client
}

func New(ctx context.Context, usr, pass string, oauth Oauth2) (*retryablehttp.Client, error) {
	state := getState(10)

	action := oauth2.SetAuthURLParam("action", "Login")
	username := oauth2.SetAuthURLParam("username", usr)
	password := oauth2.SetAuthURLParam("password", pass)

	code := oauth.AuthCodeURL(state, action, username, password)
	if code == "" {
		return nil, errors.New("authorization code is empty")
	}

	token, err := oauth.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchanging code for token: %w", err)
	}

	httpClient := oauth.Client(ctx, token)
	if httpClient == nil {
		return nil, errors.New("could not instantiate oauth http client")
	}

	return &retryablehttp.Client{
		HTTPClient:   httpClient,
		Logger:       logger,
		RetryWaitMin: retryWaitMin,
		RetryWaitMax: retryWaitMax,
		RetryMax:     retryMax,
		CheckRetry:   retryablehttp.DefaultRetryPolicy,
		Backoff:      retryablehttp.DefaultBackoff,
	}, nil
}

func getState(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}
