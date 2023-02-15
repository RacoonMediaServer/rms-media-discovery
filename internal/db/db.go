package db

import (
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/RacoonMediaServer/rms-packages/pkg/configuration"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	AccountDatabase
}

type database struct {
	conn *gorm.DB
}

func Connect(config configuration.Database) (Database, error) {
	db, err := gorm.Open(postgres.Open(config.GetConnectionString()))
	if err != nil {
		return nil, err
	}

	if err = db.AutoMigrate(&model.Account{}); err != nil {
		return nil, err
	}

	return &database{conn: db}, nil
}
