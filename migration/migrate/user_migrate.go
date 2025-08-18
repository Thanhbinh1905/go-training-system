package migrate

import (
	"fmt"

	"github.com/Thanhbinh1905/go-training-system/migration/model"
	"github.com/Thanhbinh1905/go-training-system/shared/db"
	"go.uber.org/zap"
)

func RunUserMigrations(dbURL string, log *zap.Logger) error {
	conn, err := db.Connect(dbURL, log)
	if err != nil {
		log.Fatal("failed to connect to database: ", zap.Error(err))
	}
	defer db.Close(conn, log)

	if err := conn.AutoMigrate(&model.User{}); err != nil {
		return fmt.Errorf("failed to run user migrations: %w", err)
	}

	return nil
}
