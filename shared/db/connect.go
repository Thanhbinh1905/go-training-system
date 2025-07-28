package db

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	gormlogger "gorm.io/gorm/logger"
)

func Connect(dbURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: gormlogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			gormlogger.Config{
				LogLevel: gormlogger.Silent,
			},
		),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Close(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		return
	}

	if err := sqlDB.Close(); err != nil {
	} else {
	}
}
