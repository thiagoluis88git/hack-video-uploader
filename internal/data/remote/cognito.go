package remote

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type CognitoRemoteDataSource interface {
	LoginUnknown() (string, error)
}

type CognitoRemoteDataSourceImpl struct {
	cognitoClient *cognito.CognitoIdentityProvider
	appClientID   string
}

func NewCognitoRemoteDataSource(
	region string,
	appClientId string,
) CognitoRemoteDataSource {
	config := &aws.Config{Region: aws.String(region)}
	sess, err := session.NewSession(config)
	if err != nil {
		panic(err)
	}
	client := cognito.New(sess)

	client.AdminUpdateUserAttributes(&cognito.AdminUpdateUserAttributesInput{
		UserAttributes: []*cognito.AttributeType{
			{
				Name:  aws.String("email_verified"),
				Value: aws.String("true"),
			},
		},
	})

	// teste, testee := client.CreateUserPoolRequest()
	return &CognitoRemoteDataSourceImpl{
		cognitoClient: client,
		appClientID:   appClientId,
	}
}

func (ds *CognitoRemoteDataSourceImpl) LoginUnknown() (string, error) {
	authInput := &cognito.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: aws.StringMap(map[string]string{
			"USERNAME": "unknown-user",
			"PASSWORD": "unknown-user",
		}),
		ClientId: aws.String(ds.appClientID),
	}
	result, err := ds.cognitoClient.InitiateAuth(authInput)

	if err != nil {
		return "", err
	}

	return *result.AuthenticationResult.AccessToken, nil
}
