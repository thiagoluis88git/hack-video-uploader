package usecase

import (
	"context"

	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/entity"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/repository"
)

type FinishVideoProcessUseCase interface {
	Execute(ctx context.Context, message entity.Message) error
}

type FinishVideoProcessUseCaseImpl struct {
	repo repository.UploaderRepository
}
