package postgres

import (
	stdlog "log"
	"os"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	gormlogger "gorm.io/gorm/logger"
)

func Connect(dbURL string, log *zap.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: gormlogger.New(
			stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags),
			gormlogger.Config{
				LogLevel: gormlogger.Silent,
			},
		),
	})
	if err != nil {
		log.Error("Failed to connect to database", zap.Error(err))
		return nil, err
	}

	log.Info("Database connected successfully") // ẩn password
	return db, nil
}

func Close(db *gorm.DB, log *zap.Logger) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Warn("Failed to get sqlDB from gorm", zap.Error(err))
		return
	}

	if err := sqlDB.Close(); err != nil {
		log.Warn("Failed to close database connection", zap.Error(err))
	} else {
		log.Info("Database connection closed")
	}
}
