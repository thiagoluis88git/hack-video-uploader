package main

import (
	"net/http"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/thiagoluis88git/hack-video-uploader/internal/handler"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/di"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/environment"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/httpserver"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
)

func main() {
	env := environment.LoadEnvironmentVariables()

	ds := di.ProvidesUploaderRemoteDataSource(env.Region, "bucket-aqui")
	repo := di.ProvidesUploaderRepository(ds)
	uploadFileUseCase := di.ProvidesUploadFileUseCase(repo)

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
	server.Start()
}
