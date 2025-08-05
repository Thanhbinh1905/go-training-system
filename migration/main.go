package main

import (
	"fmt"
	"log"

	"github.com/Thanhbinh1905/go-training-system/migration/config"
	"github.com/Thanhbinh1905/go-training-system/migration/migrate"
)

func RunMigrations(name, dbURL string, migrateFn func(string) error) {
	fmt.Printf("Running %s migrations...\n", name)
	if err := migrateFn(dbURL); err != nil {
		log.Fatalf("failed to run %s migrations: %v", name, err)
	}
}

func main() {
	cfg, err := config.LoadMigrationConfig()
	if err != nil {
		log.Fatal("failed to load configuration")
	}

	RunMigrations("user", cfg.UserDBURL, migrate.RunUserMigrations)
	RunMigrations("team", cfg.TeamDBURL, migrate.RunTeamMigrations)
	RunMigrations("asset", cfg.AssetDBURL, migrate.RunAssetMigrations)

	fmt.Println("All migrations completed successfully")
}
