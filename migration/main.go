package main

import (
	"fmt"
	"log"

	"github.com/Thanhbinh1905/go-training-system/migration/config"
	"github.com/Thanhbinh1905/go-training-system/migration/migrate"
	"github.com/Thanhbinh1905/go-training-system/shared/logger"
	"go.uber.org/zap"
)

func RunMigrations(logFilePath, serviceName, dbURL string, migrateFunc func(string, *zap.Logger) error) error {
	log := logger.InitLogger(logFilePath, serviceName)
	log.Info("Starting migration", zap.String("service", serviceName))

	if err := migrateFunc(dbURL, log); err != nil {
		log.Error("Migration failed", zap.Error(err))
		return err
	}

	log.Info("Migration completed successfully", zap.String("service", serviceName))
	return nil
}

func main() {
	cfg, err := config.LoadMigrationConfig()
	if err != nil {
		log.Fatal("failed to load configuration")
	}

	RunMigrations("logs/user_migration.log", "user", cfg.UserDBURL, migrate.RunUserMigrations)
	RunMigrations("logs/team_migration.log", "team", cfg.TeamDBURL, migrate.RunTeamMigrations)
	RunMigrations("logs/asset_migration.log", "asset", cfg.AssetDBURL, migrate.RunAssetMigrations)

	fmt.Println("All migrations completed successfully")
}
