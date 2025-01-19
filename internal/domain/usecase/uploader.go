package usecase

import (
	"bytes"
	"context"
	"io"

	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/entity"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/repository"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/identity"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/queue"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
)

type UploadFileUseCase interface {
	Execute(ctx context.Context, form entity.UoloaderDocumentEntity) error
}

type UploadFileUseCaseImpl struct {
	repo         repository.UploaderRepository
	id           identity.UUIDGenerator
	queueManeger queue.QueueManager
}

func NewUploadFileUseCase(
	repo repository.UploaderRepository,
	id identity.UUIDGenerator,
	queueManeger queue.QueueManager,
) UploadFileUseCase {
	return &UploadFileUseCaseImpl{
		repo:         repo,
		id:           id,
		queueManeger: queueManeger,
	}
}

func (uc *UploadFileUseCaseImpl) Execute(ctx context.Context, form entity.UoloaderDocumentEntity) error {
	buf := bytes.NewBuffer(nil)

	if _, err := io.Copy(buf, form.Data); err != nil {
		return responses.Wrap("usecase: error when copying multipart data", err)
	}

	videoID := uc.id.New()

	err := uc.repo.UploadFile(ctx, videoID, buf.Bytes(), form.Name)

	if err != nil {
		return responses.Wrap("usecase: error when uploading file", err)
	}

	err = uc.queueManeger.WriteMessage(&videoID)

	if err != nil {
		return responses.Wrap("usecase: error when writing message in queue", err)
	}

	return nil
}
