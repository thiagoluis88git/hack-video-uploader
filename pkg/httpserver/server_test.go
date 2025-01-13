package httpserver_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/httpserver"
)

func TestServer(t *testing.T) {
	t.Parallel()

	t.Run("got success when creating http server", func(t *testing.T) {
		responseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var destination any
			err := httpserver.DecodeJSONBody(w, r, &destination)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "Request body contains badly-formed JSON (at position")

			w.WriteHeader(http.StatusOK)
		})

		httpserver.New(responseHandler)
	})

	t.Run("got success when starting http server", func(t *testing.T) {
		responseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var destination any
			err := httpserver.DecodeJSONBody(w, r, &destination)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "Request body contains badly-formed JSON (at position")

			w.WriteHeader(http.StatusOK)
		})

		s := httpserver.New(responseHandler)

		err := s.Shutdown()

		assert.NoError(t, err)
	})

	t.Run("got success when notifying http server", func(t *testing.T) {
		responseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var destination any
			err := httpserver.DecodeJSONBody(w, r, &destination)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "Request body contains badly-formed JSON (at position")

			w.WriteHeader(http.StatusOK)
		})

		s := httpserver.New(responseHandler)

		s.Notify()
	})
}
