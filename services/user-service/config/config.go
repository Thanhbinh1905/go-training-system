package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Production  bool
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()

	// Validate biến bắt buộc
	requiredVars := []string{"DATABASE_URL", "JWT_SECRET"}

	for _, key := range requiredVars {
		if !viper.IsSet(key) {
			return nil, fmt.Errorf("missing required env variable: %s", key)
		}
	}

	return &Config{
		DatabaseURL: viper.GetString("DATABASE_URL"),
		JWTSecret:   viper.GetString("JWT_SECRET"),
		Production:  viper.GetBool("PRODUCTION"), // default false nếu chưa set
	}, nil
}
