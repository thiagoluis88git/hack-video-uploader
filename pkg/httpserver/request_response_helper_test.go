package httpserver_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/httpserver"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
)

func TestRequestResponseHelper(t *testing.T) {
	t.Parallel()

	t.Run("get success when calling SendResponseSuccessWithStatus", func(t *testing.T) {
		t.Parallel()

		responseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var destination any
			err := httpserver.DecodeJSONBody(w, r, &destination)

			assert.Error(t, err)
			assert.Equal(t, "Request body must not be empty", err.Error())

			httpserver.SendResponseSuccessWithStatus(w, "", http.StatusOK)
		})

		ts := httptest.NewServer(responseHandler)
		defer ts.Close()

		req, _ := http.NewRequest(http.MethodPost, ts.URL+"/mock", strings.NewReader(""))
		req.Header.Add("Content-Type", "application/json")

		response, err := ts.Client().Do(req)

		assert.NoError(t, err)
		defer response.Body.Close()
	})

	t.Run("get success when calling SendResponseSuccess", func(t *testing.T) {
		t.Parallel()

		responseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var destination any
			err := httpserver.DecodeJSONBody(w, r, &destination)

			assert.Error(t, err)
			assert.Equal(t, "Request body must not be empty", err.Error())

			httpserver.SendResponseSuccess(w, "")
		})

		ts := httptest.NewServer(responseHandler)
		defer ts.Close()

		req, _ := http.NewRequest(http.MethodPost, ts.URL+"/mock", strings.NewReader(""))
		req.Header.Add("Content-Type", "application/json")

		response, err := ts.Client().Do(req)

		assert.NoError(t, err)
		defer response.Body.Close()
	})

	t.Run("get success when calling SendResponseNoContentSuccess", func(t *testing.T) {
		t.Parallel()

		responseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var destination any
			err := httpserver.DecodeJSONBody(w, r, &destination)

			assert.Error(t, err)
			assert.Equal(t, "Request body must not be empty", err.Error())

			httpserver.SendResponseNoContentSuccess(w)
		})

		ts := httptest.NewServer(responseHandler)
		defer ts.Close()

		req, _ := http.NewRequest(http.MethodPost, ts.URL+"/mock", strings.NewReader(""))
		req.Header.Add("Content-Type", "application/json")

		response, err := ts.Client().Do(req)

		assert.NoError(t, err)
		defer response.Body.Close()
	})

	t.Run("get success when calling SendResponseError", func(t *testing.T) {
		t.Parallel()

		responseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var destination any
			err := httpserver.DecodeJSONBody(w, r, &destination)

			assert.Error(t, err)
			assert.Equal(t, "Request body must not be empty", err.Error())

			httpserver.SendResponseError(w, errors.New("ERROR"))
		})

		ts := httptest.NewServer(responseHandler)
		defer ts.Close()

		req, _ := http.NewRequest(http.MethodPost, ts.URL+"/mock", strings.NewReader(""))
		req.Header.Add("Content-Type", "application/json")

		response, err := ts.Client().Do(req)

		assert.NoError(t, err)
		defer response.Body.Close()
	})

	t.Run("get success when calling SendBadRequestError", func(t *testing.T) {
		t.Parallel()

		responseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var destination any
			err := httpserver.DecodeJSONBody(w, r, &destination)

			assert.Error(t, err)
			assert.Equal(t, "Request body must not be empty", err.Error())

			httpserver.SendBadRequestError(w, errors.New("ERROR"))
		})

		ts := httptest.NewServer(responseHandler)
		defer ts.Close()

		req, _ := http.NewRequest(http.MethodPost, ts.URL+"/mock", strings.NewReader(""))
		req.Header.Add("Content-Type", "application/json")

		response, err := ts.Client().Do(req)

		assert.NoError(t, err)
		defer response.Body.Close()
	})

	t.Run("get StatusInternalServerError when calling GetStatusCodeFromError", func(t *testing.T) {
		t.Parallel()

		err := &responses.NetworkError{}
		status := httpserver.GetStatusCodeFromError(err)

		assert.Equal(t, http.StatusInternalServerError, status)
	})

	t.Run("get business status code when calling GetStatusCodeFromError", func(t *testing.T) {
		t.Parallel()

		err := &responses.BusinessResponse{
			StatusCode: 422,
		}

		status := httpserver.GetStatusCodeFromError(err)

		assert.Equal(t, http.StatusUnprocessableEntity, status)
	})

	t.Run("got error when passing empty json", func(t *testing.T) {
		t.Parallel()

		responseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var destination any
			err := httpserver.DecodeJSONBody(w, r, &destination)

			assert.Error(t, err)
			assert.Equal(t, "Request body must not be empty", err.Error())

			w.WriteHeader(http.StatusOK)
		})

		ts := httptest.NewServer(responseHandler)
		defer ts.Close()

		req, _ := http.NewRequest(http.MethodPost, ts.URL+"/mock", strings.NewReader(""))
		req.Header.Add("Content-Type", "application/json")

		response, err := ts.Client().Do(req)

		assert.NoError(t, err)
		defer response.Body.Close()
	})

	t.Run("got error when not passing Content-Type", func(t *testing.T) {
		t.Parallel()

		responseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var destination any
			err := httpserver.DecodeJSONBody(w, r, &destination)

			assert.Error(t, err)
			assert.Equal(t, "Content-Type header is not application/json", err.Error())

			w.WriteHeader(http.StatusOK)
		})

		ts := httptest.NewServer(responseHandler)
		defer ts.Close()

		req, _ := http.NewRequest(http.MethodPost, ts.URL+"/mock", strings.NewReader(""))

		response, err := ts.Client().Do(req)

		assert.NoError(t, err)
		defer response.Body.Close()
	})

	t.Run("got error when passing wrong Content-Type", func(t *testing.T) {
		t.Parallel()

		responseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var destination any
			err := httpserver.DecodeJSONBody(w, r, &destination)

			assert.Error(t, err)
			assert.Equal(t, "Content-Type header is not application/json", err.Error())

			w.WriteHeader(http.StatusOK)
		})

		ts := httptest.NewServer(responseHandler)
		defer ts.Close()

		req, _ := http.NewRequest(http.MethodPost, ts.URL+"/mock", strings.NewReader(""))
		req.Header.Add("Content-Type", "asdsad")

		response, err := ts.Client().Do(req)

		assert.NoError(t, err)
		defer response.Body.Close()
	})

	t.Run("got error when passing empty Content-Type", func(t *testing.T) {
		t.Parallel()

		responseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var destination any
			err := httpserver.DecodeJSONBody(w, r, &destination)

			assert.Error(t, err)
			assert.Equal(t, "Content-Type header is not application/json", err.Error())

			w.WriteHeader(http.StatusOK)
		})

		ts := httptest.NewServer(responseHandler)
		defer ts.Close()

		req, _ := http.NewRequest(http.MethodPost, ts.URL+"/mock", strings.NewReader(""))
		req.Header.Add("Content-Type", "")

		response, err := ts.Client().Do(req)

		assert.NoError(t, err)
		defer response.Body.Close()
	})

	t.Run("got error when passing badly-formed json", func(t *testing.T) {
		t.Parallel()

		responseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var destination any
			err := httpserver.DecodeJSONBody(w, r, &destination)

			assert.Error(t, err)
			assert.Equal(t, "Request body contains badly-formed JSON", err.Error())

			w.WriteHeader(http.StatusOK)
		})

		ts := httptest.NewServer(responseHandler)
		defer ts.Close()

		req, _ := http.NewRequest(http.MethodPost, ts.URL+"/mock", strings.NewReader(`{"teste": "ttt}`))
		req.Header.Add("Content-Type", "application/json")

		response, err := ts.Client().Do(req)

		assert.NoError(t, err)
		defer response.Body.Close()
	})

	t.Run("got error when passing badly-formed at position json", func(t *testing.T) {
		t.Parallel()

		responseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var destination any
			err := httpserver.DecodeJSONBody(w, r, &destination)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "Request body contains badly-formed JSON (at position")

			w.WriteHeader(http.StatusOK)
		})

		ts := httptest.NewServer(responseHandler)
		defer ts.Close()

		req, _ := http.NewRequest(http.MethodPost, ts.URL+"/mock", strings.NewReader(`{"teste": "teste"{}}`))
		req.Header.Add("Content-Type", "application/json")

		response, err := ts.Client().Do(req)

		assert.NoError(t, err)
		defer response.Body.Close()
	})

	// t.Run("got error when passing unknown fields", func(t *testing.T) {
	// 	t.Parallel()

	// 	responseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		var destination dto.Token
	// 		err := httpserver.DecodeJSONBody(w, r, &destination)

	// 		assert.Error(t, err)
	// 		assert.Contains(t, err.Error(), "Request body contains unknown field")

	// 		w.WriteHeader(http.StatusOK)
	// 	})

	// 	ts := httptest.NewServer(responseHandler)
	// 	defer ts.Close()

	// 	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/mock", strings.NewReader(`{"teste": "teste"}`))
	// 	req.Header.Add("Content-Type", "application/json")

	// 	response, err := ts.Client().Do(req)

	// 	assert.NoError(t, err)
	// 	defer response.Body.Close()
	// })

	// t.Run("got error when passing data without required fields", func(t *testing.T) {
	// 	t.Parallel()

	// 	responseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		var destination dto.Customer
	// 		err := httpserver.DecodeJSONBody(w, r, &destination)

	// 		assert.Error(t, err)
	// 		assert.Contains(t, err.Error(), "Error JSON required fields")

	// 		w.WriteHeader(http.StatusOK)
	// 	})

	// 	ts := httptest.NewServer(responseHandler)
	// 	defer ts.Close()

	// 	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/mock", strings.NewReader(`{"name": "Hamburguer"}`))
	// 	req.Header.Add("Content-Type", "application/json")

	// 	response, err := ts.Client().Do(req)

	// 	assert.NoError(t, err)
	// 	defer response.Body.Close()
	// })

	// t.Run("got success when passing valid path param", func(t *testing.T) {
	// 	t.Parallel()

	// 	responseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		id, err := httpserver.GetPathParamFromRequest(r, "id")
	// 		assert.NoError(t, err)
	// 		assert.NotEmpty(t, id)
	// 		assert.Equal(t, "123", id)

	// 		w.WriteHeader(http.StatusOK)
	// 	})

	// 	router := chi.NewRouter()

	// 	router.Get("/product/{id}", responseHandler)

	// 	ts := httptest.NewServer(router)

	// 	defer ts.Close()

	// 	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/product/123", nil)
	// 	req.Header.Add("Content-Type", "application/json")

	// 	response, err := ts.Client().Do(req)

	// 	assert.NoError(t, err)
	// 	defer response.Body.Close()
	// })

	// t.Run("got error when passing invalid path param", func(t *testing.T) {
	// 	t.Parallel()

	// 	responseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		id, err := httpserver.GetPathParamFromRequest(r, "idinvalid")
	// 		assert.Error(t, err)
	// 		assert.Empty(t, id)

	// 		var netError *responses.NetworkError
	// 		isNestError := errors.As(err, &netError)

	// 		assert.True(t, isNestError)
	// 		assert.Equal(t, http.StatusUnprocessableEntity, netError.Code)
	// 		assert.Equal(t, "idinvalid param not found in path", netError.Message)

	// 		w.WriteHeader(http.StatusUnprocessableEntity)
	// 	})

	// 	router := chi.NewRouter()

	// 	router.Get("/product/{id}", responseHandler)

	// 	ts := httptest.NewServer(router)

	// 	defer ts.Close()

	// 	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/product/123", nil)
	// 	req.Header.Add("Content-Type", "application/json")

	// 	response, err := ts.Client().Do(req)

	// 	assert.NoError(t, err)
	// 	defer response.Body.Close()
	// })
}
