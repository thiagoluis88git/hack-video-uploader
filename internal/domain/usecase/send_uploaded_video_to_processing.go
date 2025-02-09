package usecase

import (
	"context"

	"github.com/thiagoluis88git/hack-video-uploader/pkg/queue"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
)

type SendUploadedVideoForProcessingUseCase interface {
	Execute(ctx context.Context, cpf string) error
}

type SendUploadedVideoForProcessingUseCaseImpl struct {
	queueManager queue.QueueManager
}

func NewSendUploadedVideoForProcessingUseCase(
	queueManager queue.QueueManager,
) SendUploadedVideoForProcessingUseCase {
	return &SendUploadedVideoForProcessingUseCaseImpl{
		queueManager: queueManager,
	}
}

func (uc *SendUploadedVideoForProcessingUseCaseImpl) Execute(ctx context.Context, trackingID string) error {
	err := uc.queueManager.WriteMessage(&trackingID)

	if err != nil {
		return responses.Wrap("usecase: error when writing to input queue", err)
	}

	return nil
}
