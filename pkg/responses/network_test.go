package responses_test

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
)

func TestNetworkResponse(t *testing.T) {
	t.Parallel()

	t.Run("got InternalServerError error with urlError Error when calling GetNetworkError", func(t *testing.T) {
		t.Parallel()

		err := &url.Error{
			Err: errors.New(""),
		}

		localError := responses.GetNetworkError(err)

		assert.Equal(t, http.StatusInternalServerError, localError.Code)
	})

	t.Run("got StatusConflict error with Cognito Error when calling GetCognitoError", func(t *testing.T) {
		t.Parallel()

		err := errors.New("UsernameExistsException")

		localError := responses.GetCognitoError(err)

		assert.Equal(t, http.StatusConflict, localError.Code)
	})

	t.Run("got StatusInternalServerError error with Cognito Error when calling GetCognitoError", func(t *testing.T) {
		t.Parallel()

		err := errors.New("StatusInternalServerError")

		localError := responses.GetCognitoError(err)

		assert.Equal(t, http.StatusInternalServerError, localError.Code)
	})

	t.Run("got Success Status Code for responses when calling IsNetworkResponseOK", func(t *testing.T) {
		t.Parallel()

		response := &http.Response{
			StatusCode: 204,
		}

		bodyMessage := "ok"

		err := responses.IsNetworkResponseOk(response, bodyMessage)

		assert.NoError(t, err)
	})

	t.Run("got StatusCode different than Success for responses when calling IsNetworkResponseOK", func(t *testing.T) {
		t.Parallel()

		response := &http.Response{
			StatusCode: 400,
		}

		bodyMessage := "ok"

		err := responses.IsNetworkResponseOk(response, bodyMessage)

		assert.Error(t, err)
	})
}
