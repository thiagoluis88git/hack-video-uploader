package di

import (
	"fmt"

	"github.com/thiagoluis88git/hack-video-uploader/internal/data/remote"
	dataRepo "github.com/thiagoluis88git/hack-video-uploader/internal/data/repository"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/repository"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/usecase"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/storage"
)

func ProvidesUploaderRemoteDataSource(region string, bucket string) remote.UploaderRemoteDataSource {
	s3, err := storage.NewAWSS3Session(region)

	if err != nil {
		panic(fmt.Sprintf("error when getting S3 sessiogn"))
	}

	return remote.NewUploaderRemoteDataSource(s3, bucket)
}

func ProvidesUploaderRepository(ds remote.UploaderRemoteDataSource) repository.UploaderRepository {
	return dataRepo.NewUploaderRepository(ds)
}

func ProvidesUploadFileUseCase(repo repository.UploaderRepository) usecase.UploadFileUseCase {
	return usecase.NewUploadFileUseCase(repo)
}
