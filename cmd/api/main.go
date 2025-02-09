package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/thiagoluis88git/hack-video-uploader/internal/handler"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/di"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/environment"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/httpserver"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/queue"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
)

func main() {
	env := environment.LoadEnvironmentVariables()

	queueManager := queue.ConfigQueueManager(env)

	uploaderDS := di.ProvidesUploaderRemoteDataSource(env)
	cognitoDS := di.ProvidesCognitoRemoteDataSource(env)

	db := di.ProvidesDatabase(env)
	trackingLocal := di.ProvidesUploaderLocalDataSource(db)
	uploaderRepo := di.ProvidesUploaderRepository(uploaderDS, trackingLocal)
	customerRepo := di.ProvidesCustomerRepository(cognitoDS, db)
	uploadFileUseCase := di.ProvidesUploadFileUseCase(uploaderRepo, queueManager)
	finishVideoProcessUseCase := di.ProvidesFinishVideoProcessUseCase(uploaderRepo, queueManager)
	finishVideoProcessWithErrorUseCase := di.ProvidesFinishVideoProcessWithErrorUseCase(uploaderRepo, queueManager)
	presignUseCase := di.ProvidesPresignForUploadUseCase(uploaderRepo)
	getTrackingUseCase := di.ProvidesGetTrackingsUseCase(uploaderRepo)
	createUserUseCase := di.ProvidesCreateUseUseCase(customerRepo)
	loginUseCase := di.ProvidesLoginUseCase(customerRepo)

	// Config API. Must be async
	router := chi.NewRouter()
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.RealIP)
	router.Use(chiMiddleware.Recoverer)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		httpserver.SendResponseSuccess(w, &responses.BusinessResponse{
			StatusCode: 200,
			Message:    "ok",
		})
	})

	router.Post("/auth/login", handler.LoginCustomerHandler(loginUseCase))
	router.Post("/auth/signup", handler.CreateUserHandler(createUserUseCase))
	router.Get("/api/upload/presign/{cpf}", handler.GetPresignURLForUpload(presignUseCase))
	router.Post("/api/upload", handler.UploadHandler(uploadFileUseCase))
	router.Get("/api/trackings/{cpf}", handler.GetTrackingsHandler(getTrackingUseCase))

	server := httpserver.New(router)
	go server.Start()

	// Config SQS.
	chnMessages := make(chan *types.Message)
	chnErrorMessages := make(chan *types.Message)

	go queueManager.PollMessages(chnMessages)
	go queueManager.PollErrorMessages(chnErrorMessages)

	for {
		select {
		case message := <-chnMessages:
			if message == nil {
				return
			}

			log.Println("main: finishing video process")

			err := finishVideoProcessUseCase.Execute(context.Background(), message)

			if err != nil {
				log.Printf("main: error when finishing video process: %v", err.Error())
			}
		case errorMessage := <-chnErrorMessages:
			if errorMessage == nil {
				return
			}

			log.Println("main: finishing with error video process")

			err := finishVideoProcessWithErrorUseCase.Execute(context.Background(), errorMessage)

			if err != nil {
				log.Printf("main: error when finishing with error video process: %v", err.Error())
			}
		}
	}
}
