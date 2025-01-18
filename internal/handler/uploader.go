package handler

import (
	"fmt"
	"net/http"

	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/entity"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/usecase"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/httpserver"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
)

func UploadHandler(
	uploadFileUseCase usecase.UploadFileUseCase,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		documentName := r.PostFormValue("documentName")
		print(documentName)

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
