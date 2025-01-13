package httpserver

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
)

func NewHTTPClient() *http.Client {
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 20 * time.Second,
	}

	return &client
}

func DoRequest[T any](
	ctx context.Context,
	client *http.Client,
	endpoint string,
	token *string,
	formData io.Reader,
	method string,
	dataResponse T,
) (T, error) {
	var empty T

	req, err := http.NewRequestWithContext(ctx, method, endpoint, formData)

	if err != nil {
		return empty, &responses.NetworkError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	req.Header.Set("Content-Type", "application/json")

	if token != nil {
		req.Header.Set("Authorization", *token)
	}

	response, err := client.Do(req)

	if err != nil {
		return empty, responses.GetNetworkError(err.(*url.Error))
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return empty, &responses.NetworkError{
			Code:    http.StatusUnprocessableEntity,
			Message: err.Error(),
		}
	}

	bodyMessage := string(body)
	err = responses.IsNetworkResponseOk(response, bodyMessage)

	if err != nil {
		return empty, err
	}

	err = json.Unmarshal(body, &dataResponse)

	if err != nil {
		return empty, &responses.NetworkError{
			Code:    http.StatusUnprocessableEntity,
			Message: err.Error(),
		}
	}

	return dataResponse, nil
}
