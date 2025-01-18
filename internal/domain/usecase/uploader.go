package usecase

import (
	"bytes"
	"context"
	"io"

	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/entity"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/repository"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/identity"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
)

type UploadFileUseCase interface {
	Execute(ctx context.Context, form entity.UoloaderDocumentEntity) error
}

type UploadFileUseCaseImpl struct {
	repo repository.UploaderRepository
	id   identity.UUIDGenerator
}

func NewUploadFileUseCase(
	repo repository.UploaderRepository,
	id identity.UUIDGenerator,
) UploadFileUseCase {
	return &UploadFileUseCaseImpl{
		repo: repo,
		id:   id,
	}
}

func (uc *UploadFileUseCaseImpl) Execute(ctx context.Context, form entity.UoloaderDocumentEntity) error {
	buf := bytes.NewBuffer(nil)

	if _, err := io.Copy(buf, form.Data); err != nil {
		return responses.Wrap("usecase: error when copying multipart data", err)
	}

	err := uc.repo.UploadFile(ctx, uc.id.New(), buf.Bytes(), form.Name)

	if err != nil {
		return responses.Wrap("usecase: error when uploading file", err)
	}

	return nil
}
