package local

import (
	"context"

	"github.com/thiagoluis88git/hack-video-uploader/internal/data/model"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/database"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
)

type TrackingLocalDataSource interface {
	SaveVideo(ctx context.Context, input model.Tracking) error
	FinishVideoProcess(ctx context.Context, trackingID string, zippedURL string, zippedPresignURL string) error
	GetTrackings(ctx context.Context) ([]model.Tracking, error)
}

type TrackingLocalDataSourceImpl struct {
	db *database.Database
}

func NewTrackingLocalDataSource(db *database.Database) TrackingLocalDataSource {
	return &TrackingLocalDataSourceImpl{
		db: db,
	}
}

func (ds *TrackingLocalDataSourceImpl) SaveVideo(ctx context.Context, input model.Tracking) error {
	err := ds.db.Connection.WithContext(ctx).Save(&input).Error

	if err != nil {
		return responses.LocalError{
			Code:    responses.DATABASE_ERROR,
			Message: err.Error(),
		}
	}

	return nil
}

func (ds *TrackingLocalDataSourceImpl) FinishVideoProcess(
	ctx context.Context,
	trackingID string,
	zippedURL string,
	zippedPresignURL string,
) error {
	err := ds.db.Connection.WithContext(ctx).
		Model(&model.Tracking{}).
		Where("tracking_id = ?", trackingID).
		Update("zip_url_file", zippedURL).
		Update("zip_url_file_presign", zippedPresignURL).
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

func (ds *TrackingLocalDataSourceImpl) GetTrackings(ctx context.Context) ([]model.Tracking, error) {
	var trackings []model.Tracking

	err := ds.db.Connection.WithContext(ctx).
		Find(&trackings).
		Error

	if err != nil {
		return []model.Tracking{}, responses.LocalError{
			Code:    responses.DATABASE_ERROR,
			Message: err.Error(),
		}
	}

	return trackings, nil
}
