package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL string
	Production  bool
	GRPCPort    string
	UserGRPCURL string

	// Redis
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	KAFKA_BROKER string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()

	// Validate biến bắt buộc
	requiredVars := []string{"DATABASE_URL", "GRPC_PORT", "PRODUCTION", "REDIS_ADDR", "USER_SERVICE_GRPC_PORT", "KAFKA_BROKER"}

	for _, key := range requiredVars {
		if !viper.IsSet(key) {
			return nil, fmt.Errorf("missing required env variable: %s", key)
		}
	}

	return &Config{
		DatabaseURL:   viper.GetString("DATABASE_URL"),
		Production:    viper.GetBool("PRODUCTION"),
		GRPCPort:      viper.GetString("GRPC_PORT"),
		RedisAddr:     viper.GetString("REDIS_ADDR"),
		RedisPassword: viper.GetString("REDIS_PASSWORD"), // optional
		RedisDB:       viper.GetInt("REDIS_DB"),          // optional, default = 0
		UserGRPCURL:   viper.GetString("USER_SERVICE_GRPC_PORT"),
		KAFKA_BROKER:  viper.GetString("KAFKA_BROKER"),
	}, nil
}
