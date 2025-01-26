package repository

import (
	"context"

	"github.com/thiagoluis88git/hack-video-uploader/internal/data/local"
	"github.com/thiagoluis88git/hack-video-uploader/internal/data/model"
	"github.com/thiagoluis88git/hack-video-uploader/internal/data/remote"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/entity"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/repository"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
)

type UploaderRepositoryImpl struct {
	ds    remote.UploaderRemoteDataSource
	local local.TrackingLocalDataSource
}

func NewUploaderRepository(
	ds remote.UploaderRemoteDataSource,
	local local.TrackingLocalDataSource,
) repository.UploaderRepository {
	return &UploaderRepositoryImpl{
		ds:    ds,
		local: local,
	}
}

func (repo *UploaderRepositoryImpl) UploadFile(ctx context.Context, key string, data []byte, description string) error {
	videoURL, err := repo.ds.UploadFile(ctx, key, data, description)

	if err != nil {
		return responses.Wrap("repository: error when uploading file", err)
	}

	input := model.Tracking{
		TrackingStatus: model.TrackingStatusProcessing,
		VideoURLFile:   videoURL,
		TrackingID:     key,
	}

	err = repo.local.SaveVideo(ctx, input)

	if err != nil {
		return responses.Wrap("repository: error when saving file in database", err)
	}

	return nil
}

func (repo *UploaderRepositoryImpl) PresignURL(ctx context.Context, key string) (string, error) {
	url, err := repo.ds.PresignURL(ctx, key)

	if err != nil {
		return "", responses.Wrap("repository: error when presigning zip url", err)
	}

	return url, nil
}

func (repo *UploaderRepositoryImpl) FinishVideoProcess(
	ctx context.Context,
	trackingID string,
	zippedURL string,
	zippedPresignURL string,
) error {
	err := repo.local.FinishVideoProcess(ctx, trackingID, zippedURL, zippedPresignURL)

	if err != nil {
		return responses.Wrap("repository: error when updating database to finish tracking", err)
	}

	return nil
}

func (repo *UploaderRepositoryImpl) GetTrackings(ctx context.Context) ([]entity.Tracking, error) {
	trackings := make([]entity.Tracking, 0)

	response, err := repo.local.GetTrackings(ctx)

	if err != nil {
		return trackings, responses.Wrap("repository: error when updating database to finish tracking", err)
	}

	for _, tracking := range response {
		trackings = append(trackings, entity.Tracking{
			TrackingID:     tracking.TrackingID,
			TrackingStatus: entity.TrackingStatus(tracking.TrackingStatus),
			VideoURLFile:   tracking.VideoURLFile,
			ZipURLFile:     tracking.ZipURLFile,
			CreatedAt:      tracking.CreatedAt,
			UpdatedAt:      tracking.UpdatedAt,
		})
	}

	return trackings, nil
}
