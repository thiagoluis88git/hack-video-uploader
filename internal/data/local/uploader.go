package local

import (
	"context"

	"github.com/thiagoluis88git/hack-video-uploader/internal/data/model"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/database"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
)

type UploaderLocalDataSource interface {
	SaveVideo(ctx context.Context, input model.Tracking) error
	FinishVideoProcess(ctx context.Context, trackingID string, zippedURL string) error
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

func (ds *UploaderLocalDataSourceImpl) FinishVideoProcess(ctx context.Context, trackingID string, zippedURL string) error {
	err := ds.db.Connection.WithContext(ctx).
		Model(&model.Tracking{}).
		Where("tracking_id = ?", trackingID).
		Update("zip_url_file", zippedURL).
		Update("tracking_status", model.TrackingStatusDone).
		Error

	if err != nil {
		return responses.LocalError{
			Code:    responses.DATABASE_ERROR,
			Message: err.Error(),
		}
	}

	return nil
}
