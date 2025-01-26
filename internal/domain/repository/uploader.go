package repository

import (
	"context"

	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/entity"
)

type UploaderRepository interface {
	UploadFile(ctx context.Context, key string, data []byte, description string) error
	PresignURL(ctx context.Context, key string) (string, error)
	FinishVideoProcess(ctx context.Context, trackingID string, zippedURL string) error
	GetTrackings(ctx context.Context) ([]entity.Tracking, error)
}
