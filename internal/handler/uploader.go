package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/thiagoluis88git/hack-video-uploader/internal/data/remote"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/httpserver"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/storage"
)

func UploadHandler() http.HandlerFunc {
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
		// size := handler.Size
		// contentType := handler.Header.Get("Content-Type")

		// form := entity.UoloaderDocumentEntity{
		// 	Data:        file,
		// 	Name:        fileName,
		// 	Size:        size,
		// 	ContentType: contentType,
		// }

		s3, err := storage.NewAWSS3Session("us-east-1")

		if err != nil {
			httpserver.SendBadRequestError(w, &responses.BusinessResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    fmt.Sprintf("Error while getting s3 session: %v", err.Error()),
			})

			return
		}

		buf := bytes.NewBuffer(nil)

		if _, err := io.Copy(buf, file); err != nil {
			httpserver.SendBadRequestError(w, &responses.BusinessResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    fmt.Sprintf("Error while getting s3 session: %v", err.Error()),
			})

			return
		}

		ds := remote.NewUploaderRemoteDataSource(s3, "bucket-aqui")

		err = ds.UploadFile(ctx, fmt.Sprintf("doc-%v-%v", fileName, time.Now().GoString()), buf.Bytes(), "alguma coisa aqui")

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
