package migrate

import (
	"fmt"

	"github.com/Thanhbinh1905/go-training-system/migration/model"
	"github.com/Thanhbinh1905/go-training-system/shared/db"
	"go.uber.org/zap"
)

func RunAssetMigrations(dbURL string, log *zap.Logger) error {
	conn, err := db.Connect(dbURL, log)
	if err != nil {
		log.Error("failed to connect to database", zap.Error(err))
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close(conn, log)

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
