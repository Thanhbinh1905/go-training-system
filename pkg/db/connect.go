package db

import (
	"log"
	"os"

	"github.com/Thanhbinh1905/go-training-system/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	gormlogger "gorm.io/gorm/logger"
)

func Connect(dbURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: gormlogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			gormlogger.Config{
				LogLevel: gormlogger.Silent, // ❗ hoặc logger.Silent nếu muốn tắt hoàn toàn
			},
		),
	})
	if err != nil {
		logger.Log.Error("failed to connect to database", zap.Error(err))
		return nil, err
	}

	logger.Log.Info("connected to database successfully")
	return db, nil
}

func Close(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		logger.Log.Error("failed to get sql.DB from gorm.DB", zap.Error(err))
		return
	}

	if err := sqlDB.Close(); err != nil {
		logger.Log.Error("failed to close database connection", zap.Error(err))
	} else {
		logger.Log.Info("database connection closed successfully")
	}
}
