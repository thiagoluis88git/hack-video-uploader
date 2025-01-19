package usecase

import (
	"context"

	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/entity"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/repository"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/queue"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
)

type FinishVideoProcessUseCase interface {
	Execute(ctx context.Context, message entity.Message) error
}

type FinishVideoProcessUseCaseImpl struct {
	repo         repository.UploaderRepository
	queueManager queue.QueueManager
}

func NewFinishVideoProcessUseCase(
	repo repository.UploaderRepository,
	queueManager queue.QueueManager,
) FinishVideoProcessUseCase {
	return &FinishVideoProcessUseCaseImpl{
		repo:         repo,
		queueManager: queueManager,
	}
}

func (uc *FinishVideoProcessUseCaseImpl) Execute(ctx context.Context, message entity.Message) error {
	err := uc.repo.FinishVideoProcess(ctx, message.TrackingID, message.ZippedURL)

	if err != nil {
		return responses.Wrap("usecase: error when saving file in database", err)
	}

	err = uc.queueManager.DeleteMessage(&message.ReceiptHandle)

	if err != nil {
		return responses.Wrap("usecase: error when deleting message in queue", err)
	}

	return nil
}
