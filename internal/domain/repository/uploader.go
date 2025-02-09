package repository

import (
	"context"

	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/entity"
)

type UploaderRepository interface {
	UploadFile(ctx context.Context, key string, data []byte, description string) error
	StartVideoUpload(ctx context.Context, key string, presignURL string, cpf string) (entity.Tracking, error)
	PresignURL(ctx context.Context, key string) (string, error)
	PresignForUploadVideoURL(ctx context.Context, key string) (string, error)
	FinishVideoProcess(ctx context.Context, trackingID string, zippedURL string, zippedPresignURL string) error
	FinishVideoProcessWithError(ctx context.Context, trackingID string, errorMessage string) error
	GetTrackings(ctx context.Context, cpf string) ([]entity.Tracking, error)
}
