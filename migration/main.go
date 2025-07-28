package main

import (
	"fmt"
	"log"

	"github.com/Thanhbinh1905/go-training-system/migration/config"
	"github.com/Thanhbinh1905/go-training-system/migration/migrate"
	"github.com/Thanhbinh1905/go-training-system/shared/db"
)

func main() {
	cfg, err := config.LoadMigrationConfig()
	if err != nil {
		log.Fatal("failed to load configuration")
	}

	conn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed connect to database")
	}
	defer db.Close(conn)

	if err := migrate.RunMigrations(conn); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("✅ Migration complete!")
}
