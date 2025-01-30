package remote

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/thiagoluis88git/hack-video-uploader/internal/data/model"
)

const (
	passwordSufixTemp = "12!@Az"
	passwordSufix     = "1234&$sWa"
)

type CognitoRemoteDataSource interface {
	SignUp(user *model.Customer) error
	Login(cpf string) (string, error)
}

type CognitoRemoteDataSourceImpl struct {
	cognitoClient *cognito.CognitoIdentityProvider
	appClientID   string
	userPoolID    string
	groupUser     string
}

func NewCognitoRemoteDataSource(
	region string,
	userPoolID string,
	appClientId string,
	groupUser string,
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
		userPoolID:    userPoolID,
		groupUser:     groupUser,
	}
}

func (ds *CognitoRemoteDataSourceImpl) SignUp(user *model.Customer) error {
	return ds.signUp(user.CPF, user.Name, user.Email, ds.groupUser)
}

func (ds *CognitoRemoteDataSourceImpl) signUp(cpf, name, email, groupName string) error {
	messageAction := "SUPPRESS"

	pass := fmt.Sprintf("%v%v", cpf, passwordSufixTemp)

	userCognito := &cognito.AdminCreateUserInput{
		UserPoolId:        aws.String(ds.userPoolID),
		Username:          aws.String(cpf),
		MessageAction:     &messageAction,
		TemporaryPassword: &pass,
		UserAttributes: []*cognito.AttributeType{
			{
				Name:  aws.String("name"),
				Value: aws.String(name),
			},
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
			{
				Name:  aws.String("email_verified"),
				Value: aws.String("True"),
			},
		},
	}

	_, err := ds.cognitoClient.AdminCreateUser(userCognito)

	if err != nil {
		return err
	}

	password := fmt.Sprintf("%v%v", cpf, passwordSufix)
	permanent := true

	setPasswordInput := &cognito.AdminSetUserPasswordInput{
		Password:   &password,
		UserPoolId: aws.String(ds.userPoolID),
		Username:   aws.String(cpf),
		Permanent:  &permanent,
	}

	_, errPasswd := ds.cognitoClient.AdminSetUserPassword(setPasswordInput)

	if errPasswd != nil {
		return errPasswd
	}

	addUserToGroupInput := &cognito.AdminAddUserToGroupInput{
		GroupName:  &groupName,
		UserPoolId: &ds.userPoolID,
		Username:   aws.String(cpf),
	}

	_, errGroup := ds.cognitoClient.AdminAddUserToGroup(addUserToGroupInput)

	if errGroup != nil {
		return errGroup
	}

	return nil
}

func (ds *CognitoRemoteDataSourceImpl) Login(cpf string) (string, error) {
	password := fmt.Sprintf("%v%v", cpf, passwordSufix)

	authInput := &cognito.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: aws.StringMap(map[string]string{
			"USERNAME": cpf,
			"PASSWORD": password,
		}),
		ClientId: aws.String(ds.appClientID),
	}
	result, err := ds.cognitoClient.InitiateAuth(authInput)

	if err != nil {
		return "", err
	}

	return *result.AuthenticationResult.AccessToken, nil
}
