package migrate

import (
	"fmt"

	"github.com/Thanhbinh1905/go-training-system/migration/model"
	"github.com/Thanhbinh1905/go-training-system/shared/db/postgres"
	"go.uber.org/zap"
)

func RunTeamMigrations(dbURL string, log *zap.Logger) error {
	conn, err := postgres.Connect(dbURL, log)
	if err != nil {
		log.Fatal("failed to connect to database: ", zap.Error(err))
	}
	defer postgres.Close(conn, log)

	if err := conn.AutoMigrate(
		&model.Team{},
		&model.TeamMember{},
		&model.TeamManager{},
	); err != nil {
		return fmt.Errorf("failed to run user migrations: %w", err)
	}

	return nil
}
