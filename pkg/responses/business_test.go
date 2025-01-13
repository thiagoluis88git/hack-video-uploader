package responses_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
)

func TestBusinessResponse(t *testing.T) {
	t.Parallel()

	t.Run("got BadRequest error with Network Error when calling GetResponseError", func(t *testing.T) {
		t.Parallel()

		err := &responses.NetworkError{
			Code:    http.StatusBadRequest,
			Message: "BadRequest",
		}

		businessError := responses.GetResponseError(err, "MOCK")

		assert.Equal(t, "Bad request trying to execute MOCK - BadRequest", businessError.Error())
	})

	t.Run("got Unauthorized error with Network Error when calling GetResponseError", func(t *testing.T) {
		t.Parallel()

		err := &responses.NetworkError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}

		businessError := responses.GetResponseError(err, "MOCK")

		assert.Equal(t, "Unauthorized error trying to execute MOCK - Unauthorized", businessError.Error())
	})

	t.Run("got Unauthorized error with Network Error when calling GetResponseError", func(t *testing.T) {
		t.Parallel()

		err := &responses.NetworkError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}

		businessError := responses.GetResponseError(err, "MOCK")

		assert.Equal(t, "Unauthorized error trying to execute MOCK - Unauthorized", businessError.Error())
	})

	t.Run("got Forbidden error with Network Error when calling GetResponseError", func(t *testing.T) {
		t.Parallel()

		err := &responses.NetworkError{
			Code:    http.StatusForbidden,
			Message: "Forbidden",
		}

		businessError := responses.GetResponseError(err, "MOCK")

		assert.Equal(t, "Forbiden error trying to execute MOCK - Forbidden", businessError.Error())
	})

	t.Run("got NotFound error with Network Error when calling GetResponseError", func(t *testing.T) {
		t.Parallel()

		err := &responses.NetworkError{
			Code:    http.StatusNotFound,
			Message: "NotFound",
		}

		businessError := responses.GetResponseError(err, "MOCK")

		assert.Equal(t, "Not found trying to execute MOCK - NotFound", businessError.Error())
	})

	t.Run("got Conflict error with Network Error when calling GetResponseError", func(t *testing.T) {
		t.Parallel()

		err := &responses.NetworkError{
			Code:    http.StatusConflict,
			Message: "Conflict",
		}

		businessError := responses.GetResponseError(err, "MOCK")

		assert.Equal(t, "Conflit with some data using the service MOCK - Conflict", businessError.Error())
	})

	t.Run("got UnprocessableEntity error with Network Error when calling GetResponseError", func(t *testing.T) {
		t.Parallel()

		err := &responses.NetworkError{
			Code:    http.StatusUnprocessableEntity,
			Message: "UnprocessableEntity",
		}

		businessError := responses.GetResponseError(err, "MOCK")

		assert.Equal(t, "Logic error found in service MOCK - UnprocessableEntity", businessError.Error())
	})

	t.Run("got InternalServerError error with Network Error when calling GetResponseError", func(t *testing.T) {
		t.Parallel()

		err := &responses.NetworkError{
			Code:    http.StatusInternalServerError,
			Message: "InternalServerError",
		}

		businessError := responses.GetResponseError(err, "MOCK")

		assert.Equal(t, "Unexpected internal error trying to execute service MOCK - InternalServerError", businessError.Error())
	})

	t.Run("got NotFound error with Local Error when calling GetResponseError", func(t *testing.T) {
		t.Parallel()

		err := &responses.LocalError{
			Code:    responses.NOT_FOUND_ERROR,
			Message: "NotFound",
		}

		businessError := responses.GetResponseError(err, "MOCK")

		assert.Equal(t, http.StatusNotFound, businessError.(*responses.BusinessResponse).StatusCode)
	})

	t.Run("got StatusConflict error with Local Error when calling GetResponseError", func(t *testing.T) {
		t.Parallel()

		err := &responses.LocalError{
			Code:    responses.DATABASE_CONFLICT_ERROR,
			Message: "StatusConflict",
		}

		businessError := responses.GetResponseError(err, "MOCK")

		assert.Equal(t, http.StatusConflict, businessError.(*responses.BusinessResponse).StatusCode)
	})

	t.Run("got StatusServiceUnavailable error with Local Error when calling GetResponseError", func(t *testing.T) {
		t.Parallel()

		err := &responses.LocalError{
			Code:    responses.DATABASE_ERROR,
			Message: "StatusServiceUnavailable",
		}

		businessError := responses.GetResponseError(err, "MOCK")

		assert.Equal(t, http.StatusServiceUnavailable, businessError.(*responses.BusinessResponse).StatusCode)
	})

	t.Run("got StatusUnprocessableEntity error with Local Error when calling GetResponseError", func(t *testing.T) {
		t.Parallel()

		err := &responses.LocalError{
			Code:    responses.DATABASE_CONSTRAINT_ERROR,
			Message: "StatusUnprocessableEntity",
		}

		businessError := responses.GetResponseError(err, "MOCK")

		assert.Equal(t, http.StatusUnprocessableEntity, businessError.(*responses.BusinessResponse).StatusCode)
	})

	t.Run("got StatusUnprocessableEntity error with BusinessResponse Error when calling GetResponseError", func(t *testing.T) {
		t.Parallel()

		err := &responses.BusinessResponse{
			StatusCode: 422,
			Message:    "StatusUnprocessableEntity",
		}

		businessError := responses.GetResponseError(err, "MOCK")

		assert.Equal(t, http.StatusUnprocessableEntity, businessError.(*responses.BusinessResponse).StatusCode)
	})
}
