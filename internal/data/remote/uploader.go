package remote

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
)

type UploaderRemoteDataSource interface {
	UploadFile(ctx context.Context, key string, data []byte, description string) (string, error)
	PresignURL(ctx context.Context, key string) (string, error)
	PresignForUploadVideoURL(ctx context.Context, key string) (string, error)
}

type AWSS3UploaderRemoteDataSourceImpl struct {
	session   s3iface.S3API
	bucket    string
	bucketZip string
}

func NewUploaderRemoteDataSource(session s3iface.S3API, bucket string, bucketZip string) UploaderRemoteDataSource {
	return &AWSS3UploaderRemoteDataSourceImpl{
		session:   session,
		bucket:    bucket,
		bucketZip: bucketZip,
	}
}

func (ds *AWSS3UploaderRemoteDataSourceImpl) UploadFile(ctx context.Context, key string, data []byte, description string) (string, error) {
	buffer := bytes.NewBuffer(data)

	uploader := s3manager.NewUploaderWithClient(ds.session)

	output, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket:             aws.String(ds.bucket),
		Key:                aws.String(key),
		Body:               aws.ReadSeekCloser(buffer),
		ContentDisposition: aws.String(description),
		ContentType:        aws.String(http.DetectContentType(data)),
	})

	if err != nil {
		return "", responses.Wrap("AWS S3 upload error", err)
	}

	return output.Location, nil
}

func (ds *AWSS3UploaderRemoteDataSourceImpl) PresignURL(ctx context.Context, key string) (string, error) {
	svc := ds.session.(*s3.S3)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket:                     aws.String(ds.bucketZip),
		Key:                        aws.String(key),
		ResponseContentDisposition: aws.String("attachment"),
	})

	urlStr, err := req.Presign(15 * time.Minute)

	if err != nil {
		return "", responses.Wrap("AWS S3 presign url error", err)
	}

	return urlStr, nil
}

func (ds *AWSS3UploaderRemoteDataSourceImpl) PresignForUploadVideoURL(ctx context.Context, key string) (string, error) {
	svc := ds.session.(*s3.S3)

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(ds.bucket),
		Key:    aws.String(key),
	})

	urlStr, err := req.Presign(15 * time.Minute)

	if err != nil {
		return "", responses.Wrap("AWS S3 presign for put object url error", err)
	}

	return urlStr, nil
}
