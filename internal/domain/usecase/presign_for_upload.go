package usecase

import (
	"context"

	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/entity"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/repository"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/identity"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
)

type PresignForUploadUseCase interface {
	Execute(ctx context.Context, cpf string) (entity.Tracking, error)
}

type PresignForUploadUseCaseImpl struct {
	repo repository.UploaderRepository
	id   identity.UUIDGenerator
}

func NewPresignForUploadUseCase(
	repo repository.UploaderRepository,
	id identity.UUIDGenerator,
) PresignForUploadUseCase {
	return &PresignForUploadUseCaseImpl{
		repo: repo,
		id:   id,
	}
}

func (uc *PresignForUploadUseCaseImpl) Execute(ctx context.Context, cpf string) (entity.Tracking, error) {
	trackingID := uc.id.New()

	presignURL, err := uc.repo.PresignForUploadVideoURL(ctx, trackingID)

	if err != nil {
		return entity.Tracking{}, responses.Wrap("usecase: error when presigning zip url", err)
	}

	tracking, err := uc.repo.StartVideoUpload(ctx, trackingID, presignURL, cpf)

	if err != nil {
		return entity.Tracking{}, responses.Wrap("usecase: error when starting url", err)
	}

	return tracking, nil
}
