package remote_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thiagoluis88git/hack-video-uploader/internal/data/remote"
)

func TestUploaderRemote(t *testing.T) {
	t.Parallel()

	t.Run("got success when uploading file data source", func(t *testing.T) {
		t.Parallel()

		s3 := mockS3(emptyList(), 200)

		s3Mock := S3Mock{
			S3API: s3,
			files: map[string][]byte{},
			tags:  map[string]map[string]string{},
		}

		ds := remote.NewUploaderRemoteDataSource(s3Mock, "bucket")

		videoURL, err := ds.UploadFile(context.TODO(), "bucketKey", []byte("bytes da imagem"), "descricao")

		assert.NotEmpty(t, videoURL)
		assert.NoError(t, err)
	})

	t.Run("got error when uploading file data source", func(t *testing.T) {
		t.Parallel()

		s3 := mockS3(emptyList(), 500)

		s3Mock := S3Mock{
			S3API: s3,
			files: map[string][]byte{},
			tags:  map[string]map[string]string{},
		}

		ds := remote.NewUploaderRemoteDataSource(s3Mock, "bucket")

		videoURL, err := ds.UploadFile(context.TODO(), "bucketKey", []byte("bytes da imagem"), "descricao")

		assert.Empty(t, videoURL)
		assert.Error(t, err)
	})
}
