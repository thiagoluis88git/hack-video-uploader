package repository

import (
	"context"

	"github.com/thiagoluis88git/hack-video-uploader/internal/data/local"
	"github.com/thiagoluis88git/hack-video-uploader/internal/data/model"
	"github.com/thiagoluis88git/hack-video-uploader/internal/data/remote"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/entity"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/repository"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/utils"
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

	cpf := ctx.Value(utils.CtxKeyCPF{}).(string)

	input := model.Tracking{
		TrackingStatus: model.TrackingStatusProcessing,
		VideoURLFile:   videoURL,
		TrackingID:     key,
		CPF:            cpf,
	}

	err = repo.local.SaveVideo(ctx, input)

	if err != nil {
		return responses.Wrap("repository: error when saving file in database", err)
	}

	return nil
}

// StartVideoUpload implements repository.UploaderRepository.
func (repo *UploaderRepositoryImpl) StartVideoUpload(ctx context.Context, key string, presignURL string, cpf string) (entity.Tracking, error) {
	input := model.Tracking{
		TrackingStatus: model.TrackingStatusStarted,
		VideoURLFile:   presignURL,
		TrackingID:     key,
		CPF:            cpf,
	}

	err := repo.local.SaveVideo(ctx, input)

	if err != nil {
		return entity.Tracking{}, responses.Wrap("repository: error when saving file in database", err)
	}

	return entity.Tracking{
		TrackingID:     key,
		VideoURLFile:   presignURL,
		TrackingStatus: model.TrackingStatusStarted,
	}, nil
}

func (repo *UploaderRepositoryImpl) PresignURL(ctx context.Context, key string) (string, error) {
	url, err := repo.ds.PresignURL(ctx, key)

	if err != nil {
		return "", responses.Wrap("repository: error when presigning zip url", err)
	}

	return url, nil
}

func (repo *UploaderRepositoryImpl) PresignForUploadVideoURL(ctx context.Context, key string) (string, error) {
	url, err := repo.ds.PresignForUploadVideoURL(ctx, key)

	if err != nil {
		return "", responses.Wrap("repository: error when presigning url for upload", err)
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

func (repo *UploaderRepositoryImpl) GetTrackings(ctx context.Context, cpf string) ([]entity.Tracking, error) {
	trackings := make([]entity.Tracking, 0)

	response, err := repo.local.GetTrackings(ctx, cpf)

	if err != nil {
		return trackings, responses.Wrap("repository: error when updating database to finish tracking", err)
	}

	for _, tracking := range response {
		trackings = append(trackings, entity.Tracking{
			TrackingID:     tracking.TrackingID,
			TrackingStatus: entity.TrackingStatus(tracking.TrackingStatus),
			VideoURLFile:   tracking.VideoURLFile,
			ZipURLFile:     tracking.ZipURLFilePresign,
			CreatedAt:      tracking.CreatedAt,
			UpdatedAt:      tracking.UpdatedAt,
			ErrorMessage:   tracking.ErrorMessage,
		})
	}

	return trackings, nil
}

func (repo *UploaderRepositoryImpl) FinishVideoProcessWithError(ctx context.Context, trackingID string, errorMessage string) error {
	err := repo.local.FinishVideoProcessWithError(ctx, trackingID, errorMessage)

	if err != nil {
		return responses.Wrap("repository: error when updating database to finish tracking with error", err)
	}

	return nil
}
