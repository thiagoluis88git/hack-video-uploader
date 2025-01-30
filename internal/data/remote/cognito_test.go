package remote_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thiagoluis88git/hack-video-uploader/internal/data/remote"
)

func TestCognitoRemote(t *testing.T) {
	t.Parallel()

	t.Run("got error when login cognito remote", func(t *testing.T) {
		sut := remote.NewCognitoRemoteDataSource("region", "userPoolID", "appClient", "groupUser")

		result, err := sut.Login("cpf")
		assert.Error(t, err)
		assert.Empty(t, result)
	})

}
