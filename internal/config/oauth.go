package config

import "golang.org/x/oauth2"

// Oauth defines methods to retrieve oauth configuration
type Oauth interface {
	OauthConfig() *oauth2.Config
}

// OauthConfig returns the oauth config
func (c Config) OauthConfig() *oauth2.Config {
	return c.Get("auth.config").(*oauth2.Config)
}
