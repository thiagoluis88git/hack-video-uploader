package local

import (
	"context"

	"github.com/thiagoluis88git/hack-video-uploader/pkg/database"
)

type UploaderLocalDataSource interface {
	SaveVideo(ctx context.Context, videoURL string) error
}

type UploaderLocalDataSourceImpl struct {
	db *database.Database
}

func NewUploaderLocalDataSource(db *database.Database) UploaderLocalDataSource {
	return &UploaderLocalDataSourceImpl{
		db: db,
	}
}

func (u *UploaderLocalDataSourceImpl) SaveVideo(ctx context.Context, videoURL string) error {
	panic("unimplemented")
}
