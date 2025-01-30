package usecase

import (
	"context"

	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/entity"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/repository"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
)

type GetTrackingsUseCase interface {
	Execute(ctx context.Context, cpf string) ([]entity.Tracking, error)
}

type GetTrackingsUseCaseImpl struct {
	repo repository.UploaderRepository
}

func NewGetTrackingsUseCase(repo repository.UploaderRepository) GetTrackingsUseCase {
	return &GetTrackingsUseCaseImpl{
		repo: repo,
	}
}

func (uc *GetTrackingsUseCaseImpl) Execute(ctx context.Context, cpf string) ([]entity.Tracking, error) {
	trackings, err := uc.repo.GetTrackings(ctx, cpf)

	if err != nil {
		return []entity.Tracking{}, responses.Wrap("usecase: error when getting trackings", err)
	}

	return trackings, nil
}
