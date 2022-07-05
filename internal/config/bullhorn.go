package config

// Bullhorn defines methods to retrieve bullhorn configuration
type Bullhorn interface {
	BullhornUsername() string
	BullhornPassword() string
	BullhornSubscriptionUrl() string
	BullhornEntityUrl() string
	BullhornClientID() string
	BullhornSecret() string
	BullhornAuthUrl() string
	BullhornTokenUrl() string
	BullhornRedirectUrl() string
}

// BullhornUsername returns the bullhorn username
func (c Config) BullhornUsername() string {
	return c.GetString("bullhorn.username")
}

// BullhornPassword returns the bullhorn password
func (c Config) BullhornPassword() string {
	return c.GetString("bullhorn.password")
}

// BullhornSubscriptionUrl returns the bullhorn subscription url
func (c Config) BullhornSubscriptionUrl() string {
	return c.GetString("bullhorn.subscription_url")
}

// BullhornEntityUrl returns the bullhorn entity url
func (c Config) BullhornEntityUrl() string {
	return c.GetString("bullhorn.entity_url")
}

// BullhornClientID returns the bullhorn client id
func (c Config) BullhornClientID() string {
	return c.GetString("bullhorn.client_id")
}

// BullhornSecret returns the bullhorn secret
func (c Config) BullhornSecret() string {
	return c.GetString("bullhorn.secret")
}

// BullhornAuthUrl returns the bullhorn auth url
func (c Config) BullhornAuthUrl() string {
	return c.GetString("bullhorn.auth_url")
}

// BullhornTokenUrl returns the bullhorn token url
func (c Config) BullhornTokenUrl() string {
	return c.GetString("bullhorn.token_url")
}

// BullhornRedirectUrl returns the bullhorn redirect url
func (c Config) BullhornRedirectUrl() string {
	return c.GetString("bullhorn.redirect_url")
}
