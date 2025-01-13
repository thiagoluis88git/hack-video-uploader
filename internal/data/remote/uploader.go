package remote

import (
	"bytes"
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
)

type UploaderRemoteDataSource interface {
	UploadFile(ctx context.Context, key string, data []byte, description string) error
}

type AWSS3UploaderRemoteDataSourceImpl struct {
	session s3iface.S3API
	bucket  string
}

func NewUploaderRemoteDataSource(session s3iface.S3API, bucket string) UploaderRemoteDataSource {
	return &AWSS3UploaderRemoteDataSourceImpl{
		session: session,
		bucket:  bucket,
	}
}

func (ds *AWSS3UploaderRemoteDataSourceImpl) UploadFile(ctx context.Context, key string, data []byte, description string) error {
	buffer := bytes.NewBuffer(data)

	uploader := s3manager.NewUploaderWithClient(ds.session)

	_, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket:             aws.String(ds.bucket),
		Key:                aws.String(key),
		Body:               aws.ReadSeekCloser(buffer),
		ContentDisposition: aws.String(description),
		ContentType:        aws.String(http.DetectContentType(data)),
	})

	if err != nil {
		return responses.Wrap("AWS S3 upload error", err)
	}

	return nil
}
