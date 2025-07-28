package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL string
}

func LoadMigrationConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	if !viper.IsSet("DATABASE_URL") {
		return nil, fmt.Errorf("missing required env variable: DATABASE_URL")
	}

	return &Config{
		DatabaseURL: viper.GetString("DATABASE_URL"),
	}, nil
}
