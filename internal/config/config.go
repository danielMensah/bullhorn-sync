package config

import (
	"fmt"
	"reflect"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

type Config struct {
	AuthConfig               oauth2.Config `mapstructure:"AUTH_CONFIG"`
	KafkaAddress             string        `mapstructure:"KAFKA_ADDRESS"`
	BullhornUsername         string        `mapstructure:"BULLHORN_USERNAME"`
	BullhornPassword         string        `mapstructure:"BULLHORN_PASSWORD"`
	BullhornSubscriptionUrl  string        `mapstructure:"BULLHORN_SUBSCRIPTION_URL"`
	BullhornEntityUrl        string        `mapstructure:"BULLHORN_ENTITY_URL"`
	BullhornClientID         string        `mapstructure:"BULLHORN_CLIENT_ID"`
	BullhornSecret           string        `mapstructure:"BULLHORN_SECRET"`
	BullhornAuthUrl          string        `mapstructure:"BULLHORN_AUTH_URL"`
	BullhornTokenUrl         string        `mapstructure:"BULLHORN_TOKEN_URL"`
	BullhornRedirectUrl      string        `mapstructure:"BULLHORN_REDIRECT_URL"`
	CassandraHosts           []string      `mapstructure:"CASSANDRA_HOSTS"`
	CassandraKeyspace        string        `mapstructure:"CASSANDRA_KEYSPACE"`
	CassandraUsername        string        `mapstructure:"CASSANDRA_USERNAME"`
	CassandraPassword        string        `mapstructure:"CASSANDRA_PASSWORD"`
	PostgresConnectionString string        `mapstructure:"POSTGRES_CONNECTION_STRING"`
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
	opt := viper.DecodeHook(
		yamlStringToStruct(config),
	)
	err = viper.Unmarshal(&config, opt)

	return &config, nil
}

func yamlStringToStruct(m interface{}) func(rf reflect.Kind, rt reflect.Kind, data interface{}) (interface{}, error) {
	return func(rf reflect.Kind, rt reflect.Kind, data interface{}) (interface{}, error) {
		if rf != reflect.String || rt != reflect.Struct {
			return data, nil
		}

		raw := data.(string)
		if raw == "" {
			return m, nil
		}

		return m, yaml.UnmarshalStrict([]byte(raw), &m)
	}
}
