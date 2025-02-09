package usecase

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/entity"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/repository"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/queue"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
)

type FinishVideoProcessWithErrorUseCase interface {
	Execute(ctx context.Context, chnMessage *types.Message) error
}

type FinishVideoProcessWithErrorUseCaseImpl struct {
	repo         repository.UploaderRepository
	queueManager queue.QueueManager
}

func NewFinishVideoProcessWithErrorUseCase(
	repo repository.UploaderRepository,
	queueManager queue.QueueManager,
) FinishVideoProcessWithErrorUseCase {
	return &FinishVideoProcessWithErrorUseCaseImpl{
		repo:         repo,
		queueManager: queueManager,
	}
}

func (uc *FinishVideoProcessWithErrorUseCaseImpl) Execute(ctx context.Context, chnMessage *types.Message) error {
	message := entity.ToErrorMessage(*chnMessage.Body)

	err := uc.repo.FinishVideoProcessWithError(ctx, message.TrackingID, message.Message)

	if err != nil {
		return responses.Wrap("usecase: error when saving file in database", err)
	}

	err = uc.queueManager.DeleteMessage(chnMessage.ReceiptHandle, uc.queueManager.ErrorQueueURL)

	if err != nil {
		return responses.Wrap("usecase: error when deleting message in error queue", err)
	}

	err = uc.queueManager.DeleteMessage(&message.InputReceiptHandle, uc.queueManager.InputQueueURL)

	if err != nil {
		return responses.Wrap("usecase: error when deleting message in input queue", err)
	}

	return nil
}
