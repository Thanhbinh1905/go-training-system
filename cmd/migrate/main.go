package main

import (
	"fmt"
	"log"

	"github.com/Thanhbinh1905/go-training-system/internal/config"
	"github.com/Thanhbinh1905/go-training-system/internal/migration"
	"github.com/Thanhbinh1905/go-training-system/pkg/db"
	"github.com/Thanhbinh1905/go-training-system/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadConfig()
	if cfg == nil {
		log.Fatal("failed to load configuration")
	}

	logger.InitLogger(cfg.Production)

	conn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Log.Error("failed to connect to database", zap.Error(err))
	}
	defer db.Close(conn)

	if err := migration.RunMigrations(conn); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("✅ Migration complete!")
}
