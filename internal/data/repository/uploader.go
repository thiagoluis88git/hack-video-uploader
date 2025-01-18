package repository

import (
	"context"

	"github.com/thiagoluis88git/hack-video-uploader/internal/data/local"
	"github.com/thiagoluis88git/hack-video-uploader/internal/data/model"
	"github.com/thiagoluis88git/hack-video-uploader/internal/data/remote"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/repository"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
)

type UploaderRepositoryImpl struct {
	ds    remote.UploaderRemoteDataSource
	local local.UploaderLocalDataSource
}

func NewUploaderRepository(
	ds remote.UploaderRemoteDataSource,
	local local.UploaderLocalDataSource,
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
	}

	err = repo.local.SaveVideo(ctx, input)

	if err != nil {
		return responses.Wrap("repository: error when saving file in database", err)
	}

	return nil
}
