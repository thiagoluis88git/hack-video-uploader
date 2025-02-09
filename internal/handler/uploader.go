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
		cpf, err := httpserver.GetPathParamFromRequest(r, "cpf")

		if err != nil {
			log.Print("get trackings", map[string]any{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendBadRequestError(w, err)
			return
		}

		ctx := r.Context()

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

// @Summary Update customer
// @Description Update customer
// @Tags Customer
// @Accept json
// @Produce json
// @Param id path int true "12"
// @Param product body entity.Customer true "customer"
// @Success 204
// @Failure 400 "Customer has required fields"
// @Failure 404 "Customer not found"
// @Router /api/upload/presign/{cpf} [put]
func GetPresignURLForUpload(presignUseCase usecase.PresignForUploadUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cpf, err := httpserver.GetPathParamFromRequest(r, "cpf")

		if err != nil {
			log.Print("presign url", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendBadRequestError(w, err)
			return
		}

		tracking, err := presignUseCase.Execute(r.Context(), cpf)

		if err != nil {
			log.Print("update customer", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendResponseError(w, err)
			return
		}

		httpserver.SendResponseSuccess(w, tracking)
	}
}
