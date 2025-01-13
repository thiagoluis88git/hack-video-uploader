package responses_test

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
)

func TestDatabaseResponse(t *testing.T) {
	t.Parallel()

	t.Run("got Specific error with Database Error when calling GetDatabaseError", func(t *testing.T) {
		t.Parallel()

		err := &pgconn.PgError{
			Code: "99",
		}

		localError := responses.GetDatabaseError(err)

		assert.Equal(t, 99, localError.Code)
	})

	t.Run("got Generic error with Database Error when calling GetDatabaseError", func(t *testing.T) {
		t.Parallel()

		err := &pgconn.PgError{
			Code: "x99",
		}

		localError := responses.GetDatabaseError(err)

		assert.Equal(t, responses.DATABASE_ERROR, localError.Code)
	})

	t.Run("got Connection error with Database Error when calling GetDatabaseError", func(t *testing.T) {
		t.Parallel()

		err := &pgconn.ConnectError{
			Config: &pgconn.Config{
				Host:     "host",
				Port:     uint16(5432),
				Database: "database",
				User:     "user",
				Password: "password",
			},
		}

		localError := responses.GetDatabaseError(err)

		assert.Equal(t, "service unavailable", localError.Message)
	})

	t.Run("got NotFound default error with Database Error when calling GetDatabaseError", func(t *testing.T) {
		t.Parallel()

		err := errors.New("record not found")

		localError := responses.GetDatabaseError(err)

		assert.Equal(t, responses.NOT_FOUND_ERROR, localError.Code)
	})

	t.Run("got Conflict error with code 23505 with Database Error when calling GetDatabaseError", func(t *testing.T) {
		t.Parallel()

		err := &pgconn.PgError{
			Code: "23505",
		}

		localError := responses.GetDatabaseError(err)

		assert.Equal(t, responses.DATABASE_CONFLICT_ERROR, localError.Code)
	})
}
