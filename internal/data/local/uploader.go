package local

import (
	"context"

	"github.com/thiagoluis88git/hack-video-uploader/internal/data/model"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/database"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
)

type UploaderLocalDataSource interface {
	SaveVideo(ctx context.Context, input model.Tracking) error
}

type UploaderLocalDataSourceImpl struct {
	db *database.Database
}

func NewUploaderLocalDataSource(db *database.Database) UploaderLocalDataSource {
	return &UploaderLocalDataSourceImpl{
		db: db,
	}
}

func (ds *UploaderLocalDataSourceImpl) SaveVideo(ctx context.Context, input model.Tracking) error {
	err := ds.db.Connection.WithContext(ctx).Save(&input).Error

	if err != nil {
		return responses.LocalError{
			Code:    responses.DATABASE_ERROR,
			Message: err.Error(),
		}
	}

	return nil
}
