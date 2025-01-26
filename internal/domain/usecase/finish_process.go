package usecase

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/entity"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/repository"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/queue"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
)

type FinishVideoProcessUseCase interface {
	Execute(ctx context.Context, chnMessage *types.Message) error
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

func (uc *FinishVideoProcessUseCaseImpl) Execute(ctx context.Context, chnMessage *types.Message) error {
	message := entity.ToMessage(*chnMessage.Body)

	zippedPresignURL, err := uc.repo.PresignURL(ctx, fmt.Sprintf("%v.zip", message.TrackingID))

	if err != nil {
		return responses.Wrap("usecase: error when presigning zip url", err)
	}

	//salvar tbm o message.ZippedURL
	err = uc.repo.FinishVideoProcess(ctx, message.TrackingID, message.ZippedURL, zippedPresignURL)

	if err != nil {
		return responses.Wrap("usecase: error when saving file in database", err)
	}

	err = uc.queueManager.DeleteMessage(chnMessage.ReceiptHandle)

	if err != nil {
		return responses.Wrap("usecase: error when deleting message in queue", err)
	}

	return nil
}
