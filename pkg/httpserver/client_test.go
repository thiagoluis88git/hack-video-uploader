package httpserver_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/httpserver"
)

func TestHttpClient(t *testing.T) {
	t.Parallel()

	t.Run("got success when creating http client", func(t *testing.T) {
		t.Parallel()

		client := httpserver.NewHTTPClient()

		assert.NotEmpty(t, client)
	})

	// t.Run("got success when calling DoRequest client", func(t *testing.T) {
	// 	t.Parallel()

	// 	responseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		httpserver.SendResponseSuccess(w, dto.Token{
	// 			AccessToken: "TOKEN",
	// 		})
	// 	})

	// 	ts := httptest.NewServer(responseHandler)
	// 	defer ts.Close()

	// 	response, err := httpserver.DoRequest(context.TODO(), ts.Client(), ts.URL, nil, nil, http.MethodGet, dto.Token{})

	// 	assert.NoError(t, err)
	// 	assert.NotEmpty(t, response)
	// })

	// t.Run("got error on invalid json when calling DoRequest client", func(t *testing.T) {
	// 	t.Parallel()

	// 	responseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		w.WriteHeader(http.StatusOK)
	// 	})

	// 	ts := httptest.NewServer(responseHandler)
	// 	defer ts.Close()

	// 	response, err := httpserver.DoRequest(context.TODO(), ts.Client(), ts.URL, nil, nil, http.MethodGet, dto.Token{})

	// 	assert.Error(t, err)
	// 	assert.Empty(t, response)
	// })

	// t.Run("got error when calling DoRequest client", func(t *testing.T) {
	// 	t.Parallel()

	// 	client := httpserver.NewHTTPClient()

	// 	assert.NotEmpty(t, client)

	// 	response, err := httpserver.DoRequest(context.TODO(), client, "http://localhost", nil, nil, http.MethodGet, dto.Token{})

	// 	assert.Error(t, err)
	// 	assert.Empty(t, response)
	// })
}
