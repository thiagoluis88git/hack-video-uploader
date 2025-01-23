package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/go-chi/chi"
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

	ds := di.ProvidesUploaderRemoteDataSource(env.Region, env.S3Bucket)
	local := di.ProvidesUploaderLocalDataSource(env)
	repo := di.ProvidesUploaderRepository(ds, local)
	uploadFileUseCase := di.ProvidesUploadFileUseCase(repo, queueManager)
	finishVideoProcessUseCase := di.ProvidesFinishVideoProcessUseCase(repo, queueManager)

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

	router.Post("/api/upload", handler.UploadHandler(uploadFileUseCase))

	server := httpserver.New(router)
	go server.Start()

	// Config SQS.
	chnMessages := make(chan *types.Message)

	go queueManager.PollMessages(chnMessages)

	for message := range chnMessages {
		if message == nil {
			return
		}

		err := finishVideoProcessUseCase.Execute(context.Background(), message)

		if err != nil {
			log.Printf("main: error when finishing video process: %v", err.Error())
		}
	}
}
