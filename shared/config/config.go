package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`
	Port        string `mapstructure:"PORT"`
	Production  bool   `mapstructure:"PRODUCTION"`
}

func LoadConfig() *Config {
	godotenv.Load()

	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDB := os.Getenv("POSTGRES_DB")
	if postgresHost == "" || postgresPort == "" || postgresUser == "" || postgresPassword == "" || postgresDB == "" {
		log.Fatal("Missing required environment variables for database connection")
		return nil
	}
	databaseUrl := "postgres://" + postgresUser + ":" + postgresPassword + "@" + postgresHost + ":" + postgresPort + "/" + postgresDB + "?sslmode=disable"
	jwt := os.Getenv("JWT_SECRET")
	port := os.Getenv("PORT")
	production := os.Getenv("PRODUCTION") == "true"

	if jwt == "" || port == "" {
		log.Fatal("Missing required environment variables for JWT secret or port")
		return nil
	}

	return &Config{
		DatabaseURL: databaseUrl,
		JWTSecret:   jwt,
		Port:        port,
		Production:  production,
	}
}
