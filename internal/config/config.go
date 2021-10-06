package config

import (
	"fmt"
	"os"

	"github.com/danielMensah/bullhorn-sync-poc/internal/bullhorn"
	"golang.org/x/oauth2"
)

var (
	envBullhornUsername        = "BULLHORN_USERNAME"
	envBullhornPassword        = "BULLHORN_PASSWORD"
	envBullhornSubscriptionUrl = "BULLHORN_PASSWORD"
	envBullhornEntityUrl       = "BULLHORN_PASSWORD"
	envBullhornClientID        = "BULLHORN_CLIENT_ID"
	envBullhornSecret          = "BULLHORN_SECRET"
	envBullhornAuthUrl         = "BULLHORN_AUTH_URL"
	envBullhornTokenUrl        = "BULLHORN_TOKEN_URL"
	envBullhornRedirectUrl     = "BULLHORN_REDIRECT_URL"

	envAWSRegion = ""
	envCronSpec  = "CRON_SPEC"

	errMissingEnv = "missing environment variable"
)

type Config struct {
	BhConfig     bullhorn.Config
	Oauth2Config *oauth2.Config
	Region       string
	CronSpec     string
}

func New() (Config, error) {
	bullhornUsername, exists := os.LookupEnv(envBullhornUsername)
	if !exists {
		return Config{}, fmt.Errorf("%s: %s", errMissingEnv, envBullhornUsername)
	}
	bullhornPassword, exists := os.LookupEnv(envBullhornPassword)
	if !exists {
		return Config{}, fmt.Errorf("%s: %s", errMissingEnv, envBullhornPassword)
	}
	bullhornSubscriptionUrl, exists := os.LookupEnv(envBullhornSubscriptionUrl)
	if !exists {
		return Config{}, fmt.Errorf("%s: %s", errMissingEnv, envBullhornSubscriptionUrl)
	}
	bullhornEntityUrl, exists := os.LookupEnv(envBullhornEntityUrl)
	if !exists {
		return Config{}, fmt.Errorf("%s: %s", errMissingEnv, envBullhornEntityUrl)
	}
	bullhornClientID, exists := os.LookupEnv(envBullhornClientID)
	if !exists {
		return Config{}, fmt.Errorf("%s: %s", errMissingEnv, envBullhornClientID)
	}
	bullhornSecret, exists := os.LookupEnv(envBullhornSecret)
	if !exists {
		return Config{}, fmt.Errorf("%s: %s", errMissingEnv, envBullhornSecret)
	}
	bullhornAuthUrl, exists := os.LookupEnv(envBullhornAuthUrl)
	if !exists {
		return Config{}, fmt.Errorf("%s: %s", errMissingEnv, envBullhornAuthUrl)
	}
	bullhornTokenUrl, exists := os.LookupEnv(envBullhornTokenUrl)
	if !exists {
		return Config{}, fmt.Errorf("%s: %s", errMissingEnv, envBullhornTokenUrl)
	}
	bullhornRedirectUrl, exists := os.LookupEnv(envBullhornRedirectUrl)
	if !exists {
		return Config{}, fmt.Errorf("%s: %s", errMissingEnv, envBullhornRedirectUrl)
	}

	spec, exists := os.LookupEnv(envCronSpec)
	if !exists {
		return Config{}, fmt.Errorf("%s: %s", errMissingEnv, envCronSpec)
	}

	region, exists := os.LookupEnv(envAWSRegion)
	if !exists {
		return Config{}, fmt.Errorf("%s %s", errMissingEnv, envAWSRegion)
	}

	return Config{
		BhConfig: bullhorn.Config{
			Username:        bullhornUsername,
			Password:        bullhornPassword,
			SubscriptionUrl: bullhornSubscriptionUrl,
			EntityUrl:       bullhornEntityUrl,
		},
		Oauth2Config: &oauth2.Config{
			ClientID:     bullhornClientID,
			ClientSecret: bullhornSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:   bullhornAuthUrl,
				TokenURL:  bullhornTokenUrl,
				AuthStyle: 0,
			},
			RedirectURL: bullhornRedirectUrl,
			Scopes:      nil,
		},
		CronSpec: spec,
		Region:   region,
	}, nil
}
