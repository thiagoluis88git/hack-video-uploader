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
	"github.com/thiagoluis88git/hack-video-uploader/pkg/identity"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/queue"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/storage"
	"gorm.io/driver/postgres"
)

func ProvidesDatabase(env environment.Environment) *database.Database {
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

	return db
}
func ProvidesUploaderRemoteDataSource(env environment.Environment) remote.UploaderRemoteDataSource {
	s3, err := storage.NewAWSS3Session(env.Region)

	if err != nil {
		panic(fmt.Sprintf("error when getting S3 session: %v", err.Error()))
	}

	return remote.NewUploaderRemoteDataSource(s3, env.S3Bucket, env.S3BucketZip)
}

func ProvidesUploaderLocalDataSource(db *database.Database) local.TrackingLocalDataSource {
	return local.NewTrackingLocalDataSource(db)
}

func ProvidesUploaderRepository(
	ds remote.UploaderRemoteDataSource,
	local local.TrackingLocalDataSource,
) repository.UploaderRepository {
	return dataRepo.NewUploaderRepository(ds, local)
}

func ProvidesCognitoRemoteDataSource(env environment.Environment) remote.CognitoRemoteDataSource {
	return remote.NewCognitoRemoteDataSource(env.Region, env.UserPoolID, env.AppClientID, env.GroupUser)
}

func ProvidesCustomerRepository(
	ds remote.CognitoRemoteDataSource,
	db *database.Database,
) repository.CustomerRepository {
	return dataRepo.NewCustomerRepository(db, ds)
}

func ProvidesUploadFileUseCase(
	repo repository.UploaderRepository,
	queueManeger queue.QueueManager,
) usecase.UploadFileUseCase {
	id := identity.NewUUIDGenerator()
	return usecase.NewUploadFileUseCase(repo, id, queueManeger)
}

func ProvidesPresignForUploadUseCase(
	repo repository.UploaderRepository,
) usecase.PresignForUploadUseCase {
	id := identity.NewUUIDGenerator()
	return usecase.NewPresignForUploadUseCase(repo, id)
}

func ProvidesSendUploadedVideoForProcessingUseCase(
	queueManeger queue.QueueManager,
) usecase.SendUploadedVideoForProcessingUseCase {
	return usecase.NewSendUploadedVideoForProcessingUseCase(queueManeger)
}

func ProvidesFinishVideoProcessUseCase(
	repo repository.UploaderRepository,
	queueManager queue.QueueManager,
) usecase.FinishVideoProcessUseCase {
	return usecase.NewFinishVideoProcessUseCase(repo, queueManager)
}

func ProvidesFinishVideoProcessWithErrorUseCase(
	repo repository.UploaderRepository,
	queueManager queue.QueueManager,
) usecase.FinishVideoProcessWithErrorUseCase {
	return usecase.NewFinishVideoProcessWithErrorUseCase(repo, queueManager)
}

func ProvidesGetTrackingsUseCase(
	repo repository.UploaderRepository,
) usecase.GetTrackingsUseCase {
	return usecase.NewGetTrackingsUseCase(repo)
}

func ProvidesCreateUseUseCase(
	repo repository.CustomerRepository,
) usecase.CreateCustomerUseCase {
	validateCPFUseCase := usecase.NewValidateCPFUseCase()
	return usecase.NewCreateCustomerUseCase(validateCPFUseCase, repo)
}

func ProvidesLoginUseCase(
	repo repository.CustomerRepository,
) usecase.LoginCustomerUseCase {
	return usecase.NewLoginCustomerUseCase(repo)
}
