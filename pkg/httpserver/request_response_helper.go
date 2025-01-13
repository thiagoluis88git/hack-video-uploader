package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/golang/gddo/httputil/header"
)

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) error {
	if r.Header.Get("Content-Type") == "" {
		msg := "Content-Type header is not application/json"
		return &responses.BusinessResponse{StatusCode: http.StatusUnsupportedMediaType, Message: msg}
	}

	value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
	if value != "application/json" {
		msg := "Content-Type header is not application/json"
		return &responses.BusinessResponse{StatusCode: http.StatusUnsupportedMediaType, Message: msg}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &responses.BusinessResponse{StatusCode: http.StatusBadRequest, Message: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			return &responses.BusinessResponse{StatusCode: http.StatusBadRequest, Message: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &responses.BusinessResponse{StatusCode: http.StatusBadRequest, Message: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &responses.BusinessResponse{StatusCode: http.StatusBadRequest, Message: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &responses.BusinessResponse{StatusCode: http.StatusBadRequest, Message: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &responses.BusinessResponse{StatusCode: http.StatusRequestEntityTooLarge, Message: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})

	if !errors.Is(err, io.EOF) {
		msg := "Request body must only contain a single JSON object"
		return &responses.BusinessResponse{StatusCode: http.StatusBadRequest, Message: msg}
	}

	validate := validator.New()
	err = validate.Struct(dst)
	if err != nil {
		msg := fmt.Sprintf("Error JSON required fields: %v", err.Error())
		return &responses.BusinessResponse{StatusCode: http.StatusBadRequest, Message: msg}
	}

	return nil
}

func SendResponseError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var br *responses.BusinessResponse

	if errors.As(err, &br) {
		w.WriteHeader(br.StatusCode)
	} else {
		br = &responses.BusinessResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Unexpected internal error",
		}

		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(br)
}

func GetStatusCodeFromError(err error) int {
	var br *responses.BusinessResponse

	if errors.As(err, &br) {
		return br.StatusCode
	}

	return http.StatusInternalServerError
}

func SendBadRequestError(w http.ResponseWriter, err error) {
	SendResponseError(w, &responses.BusinessResponse{
		StatusCode: http.StatusBadRequest,
		Message:    fmt.Sprintf("Bad request: %v", err.Error()),
	})
}

func SendResponseNoContentSuccess(w http.ResponseWriter) {
	SendResponseSuccessWithStatus(w, nil, http.StatusNoContent)
}

func SendResponseSuccess(w http.ResponseWriter, data any) {
	SendResponseSuccessWithStatus(w, data, http.StatusOK)
}

func SendResponseSuccessWithStatus(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func GetPathParamFromRequest(r *http.Request, param string) (string, error) {
	value := chi.URLParam(r, param)

	if value == "" {
		return "", &responses.NetworkError{
			Code:    http.StatusUnprocessableEntity,
			Message: fmt.Sprintf("%v param not found in path", param),
		}
	}

	return value, nil
}
