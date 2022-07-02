package config

import (
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	envVars := env{
		"AUTH_CONFIG":                "",
		"KAFKA_ADDRESS":              "",
		"BULLHORN_USERNAME":          "",
		"BULLHORN_PASSWORD":          "",
		"BULLHORN_SUBSCRIPTION_URL":  "",
		"BULLHORN_ENTITY_URL":        "",
		"BULLHORN_CLIENT_ID":         "",
		"BULLHORN_SECRET":            "",
		"BULLHORN_AUTH_URL":          "",
		"BULLHORN_TOKEN_URL":         "",
		"BULLHORN_REDIRECT_URL":      "",
		"CASSANDRA_HOSTS":            "",
		"CASSANDRA_KEYSPACE":         "",
		"CASSANDRA_USERNAME":         "",
		"CASSANDRA_PASSWORD":         "",
		"POSTGRES_CONNECTION_STRING": "",
	}
	tests := []struct {
		name    string
		want    *Config
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
