package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	RPCAddress              string   `mapstructure:"RPC_ADDRESS"`
	BullhornUsername        string   `mapstructure:"BULLHORN_USERNAME"`
	BullhornPassword        string   `mapstructure:"BULLHORN_PASSWORD"`
	BullhornSubscriptionUrl string   `mapstructure:"BULLHORN_SUBSCRIPTION_URL"`
	BullhornEntityUrl       string   `mapstructure:"BULLHORN_ENTITY_URL"`
	BullhornClientID        string   `mapstructure:"BULLHORN_CLIENT_ID"`
	BullhornSecret          string   `mapstructure:"BULLHORN_SECRET"`
	BullhornAuthUrl         string   `mapstructure:"BULLHORN_AUTH_URL"`
	BullhornTokenUrl        string   `mapstructure:"BULLHORN_TOKEN_URL"`
	BullhornRedirectUrl     string   `mapstructure:"BULLHORN_REDIRECT_URL"`
	KafkaAddress            string   `mapstructure:"KAFKA_ADDRESS"`
	CassandraHosts          []string `mapstructure:"CASSANDRA_HOSTS"`
	CassandraKeyspace       string   `mapstructure:"CASSANDRA_KEYSPACE"`
	CassandraUsername       string   `mapstructure:"CASSANDRA_USERNAME"`
	CassandraPassword       string   `mapstructure:"CASSANDRA_PASSWORD"`
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %v", err)
	}

	var config Config
	err = viper.Unmarshal(&config)

	return &config, nil
}

//Oauth2Config: &oauth2.Config{
//ClientID:     bullhornClientID,
//ClientSecret: bullhornSecret,
//Endpoint: oauth2.Endpoint{
//AuthURL:   bullhornAuthUrl,
//TokenURL:  bullhornTokenUrl,
//AuthStyle: 0,
//},
//RedirectURL: bullhornRedirectUrl,
//Scopes:      nil,
//},
