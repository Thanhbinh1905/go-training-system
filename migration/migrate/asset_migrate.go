package migrate

import (
	"fmt"
	"log"

	"github.com/Thanhbinh1905/go-training-system/migration/model"
	"github.com/Thanhbinh1905/go-training-system/shared/db"
)

func RunAssetMigrations(dbURL string) error {
	conn, err := db.Connect(dbURL)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
	defer db.Close(conn)

	if err := conn.AutoMigrate(
		&model.Folder{},
		&model.Note{},
		&model.FolderShare{},
		&model.NoteShare{},
	); err != nil {
		return fmt.Errorf("failed to run user migrations: %w", err)
	}

	return nil
}
