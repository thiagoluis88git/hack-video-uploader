package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/entity"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/usecase"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/httpserver"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/utils"
)

func UploadHandler(
	uploadFileUseCase usecase.UploadFileUseCase,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		cpf := r.FormValue("cpf")

		if cpf == "" {
			httpserver.SendBadRequestError(w, &responses.BusinessResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Error while retrieving cpf value",
			})

			return
		}

		ctx = context.WithValue(ctx, utils.CtxKeyCPF{}, cpf)

		file, handler, err := r.FormFile("file")

		if err != nil {
			httpserver.SendBadRequestError(w, &responses.BusinessResponse{
				StatusCode: http.StatusBadRequest,
				Message:    fmt.Sprintf("Error while retrieving file: %v", err.Error()),
			})

			return
		}

		defer file.Close()

		fileName := handler.Filename
		size := handler.Size
		contentType := handler.Header.Get("Content-Type")

		form := entity.UoloaderDocumentEntity{
			Data:        file,
			Name:        fileName,
			Size:        size,
			ContentType: contentType,
		}

		err = uploadFileUseCase.Execute(ctx, form)

		if err != nil {
			httpserver.SendBadRequestError(w, &responses.BusinessResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    fmt.Sprintf("Error when uploading file: %v", err.Error()),
			})

			return
		}

		httpserver.SendResponseNoContentSuccess(w)
	}
}

func GetTrackingsHandler(
	getTrackingsUseCase usecase.GetTrackingsUseCase,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		cpf, err := httpserver.GetPathParamFromRequest(r, "cpf")

		if err != nil {
			log.Print("get trackings", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendBadRequestError(w, err)
			return
		}

		trackings, err := getTrackingsUseCase.Execute(ctx, cpf)

		if err != nil {
			httpserver.SendBadRequestError(w, &responses.BusinessResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    fmt.Sprintf("Error when uploading file: %v", err.Error()),
			})

			return
		}

		httpserver.SendResponseSuccess(w, trackings)
	}
}
