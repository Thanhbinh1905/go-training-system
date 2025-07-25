package repository

import (
	"gorm.io/gorm"
)

type AssetRepository interface {
}

type assetRepository struct {
	db *gorm.DB
}
