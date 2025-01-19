package repository

import "context"

type UploaderRepository interface {
	UploadFile(ctx context.Context, key string, data []byte, description string) error
	FinishVideoProcess(ctx context.Context, trackingID string, zippedURL string) error
}
