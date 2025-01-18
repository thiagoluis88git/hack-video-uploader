package database

import (
	"github.com/thiagoluis88git/hack-video-uploader/internal/data/model"
	"gorm.io/gorm"
)

type Database struct {
	Connection *gorm.DB
}

func ConfigDatabase(dialector gorm.Dialector) (*Database, error) {
	db, err := gorm.Open(dialector, &gorm.Config{})

	if err != nil {
		return &Database{}, err
	}

	db.AutoMigrate(
		&model.Tracking{},
	)

	return &Database{
		Connection: db,
	}, nil
}
