package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Production  bool
	GRPCPort    string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()

	requiredVars := []string{"DATABASE_URL", "JWT_SECRET", "PRODUCTION", "GRPC_PORT"}

	for _, key := range requiredVars {
		if !viper.IsSet(key) {
			return nil, fmt.Errorf("missing required env variable: %s", key)
		}
	}

	return &Config{
		DatabaseURL: viper.GetString("DATABASE_URL"),
		JWTSecret:   viper.GetString("JWT_SECRET"),
		Production:  viper.GetBool("PRODUCTION"),
		GRPCPort:    viper.GetString("GRPC_PORT"),
	}, nil
}
