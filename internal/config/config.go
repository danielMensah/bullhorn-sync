package config

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Interface embeds other interfaces to provide easy access to the configs
type Interface interface {
	Bullhorn
	Kafka
	Oauth
}

// Config is a wrapper around the viper config
type Config struct {
	*viper.Viper
}

var once sync.Once

// LoadConfig loads the config from the command line flags and environment variables
func LoadConfig() (Interface, error) {
	once.Do(func() {
		initFlags()
	})

	v := viper.NewWithOptions()
	pflag.Parse()

	_ = v.BindPFlags(pflag.CommandLine)
	if envFile := v.GetString("env"); envFile != "" {
		if err := loadFromEnvFile(v, envFile); err != nil {
			return nil, fmt.Errorf("failed to load env file: %w", err)
		}
	} else {
		v.AutomaticEnv()
	}

	return Config{v}, nil
}

func initFlags() {
	pflag.String("env", "", "load environment variables from file")
}

func loadFromEnvFile(v *viper.Viper, path string) error {
	v.SetConfigFile(path)
	envType := strings.TrimPrefix(filepath.Ext(path), ".")
	v.SetConfigType(envType)

	return v.ReadInConfig()
}
