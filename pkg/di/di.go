package di

import (
	"fmt"

	"github.com/thiagoluis88git/hack-video-uploader/internal/data/local"
	"github.com/thiagoluis88git/hack-video-uploader/internal/data/remote"
	dataRepo "github.com/thiagoluis88git/hack-video-uploader/internal/data/repository"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/repository"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/usecase"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/database"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/environment"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/storage"
	"gorm.io/driver/postgres"
)

func ProvidesUploaderRemoteDataSource(region string, bucket string) remote.UploaderRemoteDataSource {
	s3, err := storage.NewAWSS3Session(region)

	if err != nil {
		panic(fmt.Sprintf("error when getting S3 session: %v", err.Error()))
	}

	return remote.NewUploaderRemoteDataSource(s3, bucket)
}

func ProvidesUploaderLocalDataSource(env environment.Environment) local.UploaderLocalDataSource {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v",
		env.DBHost,
		env.DBUser,
		env.DBPassword,
		env.DBName,
		env.DBPort,
	)

	db, err := database.ConfigDatabase(postgres.Open(dsn))

	if err != nil {
		panic(fmt.Sprintf("error when starting PostgreSQL: %v", err.Error()))
	}

	return local.NewUploaderLocalDataSource(db)
}

func ProvidesUploaderRepository(
	ds remote.UploaderRemoteDataSource,
	local local.UploaderLocalDataSource,
) repository.UploaderRepository {
	return dataRepo.NewUploaderRepository(ds, local)
}

func ProvidesUploadFileUseCase(repo repository.UploaderRepository) usecase.UploadFileUseCase {
	return usecase.NewUploadFileUseCase(repo)
}
