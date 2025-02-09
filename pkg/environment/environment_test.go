package environment_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/environment"
)

func setup() {
	os.Setenv(environment.DBHost, "DBHost")
	os.Setenv(environment.DBPassword, "DBPassword")
	os.Setenv(environment.DBName, "DBName")
	os.Setenv(environment.DBPort, "DBPort")
	os.Setenv(environment.DBUser, "DBUser")
	os.Setenv(environment.Region, "Region")
	os.Setenv(environment.UserPoolID, "UserPoolID")
	os.Setenv(environment.AppClientID, "AppClientID")
	os.Setenv(environment.GroupUser, "GroupUser")
	os.Setenv(environment.S3Bucket, "Bucket")
	os.Setenv(environment.S3BucketZip, "BucketZip")
	os.Setenv(environment.VideoProcessedOutpuQueue, "OutputQueue")
	os.Setenv(environment.VideoProcessingInputQueue, "InputtQueue")
	os.Setenv(environment.VideoProcessedErrorQueue, "ErrorQueue")
}

func TestEnvironment(t *testing.T) {
	t.Parallel()
	setup()

	t.Run("got success when initializing environment", func(t *testing.T) {
		env := environment.LoadEnvironmentVariables()

		assert.Equal(t, "DBHost", env.DBHost)
		assert.Equal(t, "DBPassword", env.DBPassword)
		assert.Equal(t, "DBPort", env.DBPort)
		assert.Equal(t, "DBName", env.DBName)
		assert.Equal(t, "DBUser", env.DBUser)
		assert.Equal(t, "Region", env.Region)
	})
}
