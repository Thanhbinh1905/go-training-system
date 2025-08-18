package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	UserDBURL  string `mapstructure:"USER_DATABASE_URL"`
	TeamDBURL  string `mapstructure:"TEAM_DATABASE_URL"`
	AssetDBURL string `mapstructure:"ASSET_DATABASE_URL"`
}

func LoadMigrationConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	requiredVars := []string{"USER_DATABASE_URL", "TEAM_DATABASE_URL", "ASSET_DATABASE_URL"}
	for _, key := range requiredVars {
		if !viper.IsSet(key) {
			return nil, fmt.Errorf("missing required env variable: %s", key)
		}
	}

	return &Config{
		UserDBURL:  viper.GetString("USER_DATABASE_URL"),
		TeamDBURL:  viper.GetString("TEAM_DATABASE_URL"),
		AssetDBURL: viper.GetString("ASSET_DATABASE_URL"),
	}, nil
}
