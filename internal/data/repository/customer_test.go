package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/thiagoluis88git/hack-video-uploader/internal/data/model"
	"github.com/thiagoluis88git/hack-video-uploader/internal/data/repository"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/entity"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
)

func TestCustomerRepository(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (suite *RepositoryTestSuite) TestCreateCustomerWithSuccess() {
	// ensure that the postgres database is empty
	var customers []model.Customer
	result := suite.db.Connection.Find(&customers)
	suite.NoError(result.Error)
	suite.Empty(customers)

	mockCognito := new(MockCognitoRemoteDataSource)
	repo := repository.NewCustomerRepository(suite.db, mockCognito)

	newCustomer := entity.Customer{
		Name:  "Teste",
		CPF:   "12312312312",
		Email: "teste@teste.com",
	}

	newCustomerModel := &model.Customer{
		Name:  "Teste",
		CPF:   "12312312312",
		Email: "teste@teste.com",
	}

	mockCognito.On("SignUp", newCustomerModel).Return(nil)

	newId, err := repo.CreateCustomer(suite.ctx, newCustomer)

	suite.NoError(err)
	suite.Equal(uint(1), newId)
}

func (suite *RepositoryTestSuite) TestCreateCustomerWithConflictError() {
	// ensure that the postgres database is empty
	var customers []model.Customer
	result := suite.db.Connection.Find(&customers)
	suite.NoError(result.Error)
	suite.Empty(customers)

	mockCognito := new(MockCognitoRemoteDataSource)
	repo := repository.NewCustomerRepository(suite.db, mockCognito)

	newCustomer := entity.Customer{
		Name:  "Teste",
		CPF:   "12312312312",
		Email: "teste@teste.com",
	}

	newCustomerModel := &model.Customer{
		Name:  "Teste",
		CPF:   "12312312312",
		Email: "teste@teste.com",
	}

	mockCognito.On("SignUp", newCustomerModel).Return(nil)

	newId, err := repo.CreateCustomer(suite.ctx, newCustomer)

	suite.NoError(err)
	suite.Equal(uint(1), newId)

	newIdError, err := repo.CreateCustomer(suite.ctx, newCustomer)

	suite.Error(err)
	suite.Equal(uint(0), newIdError)
}

func (suite *RepositoryTestSuite) TestCreateCustomerWithSignupError() {
	// ensure that the postgres database is empty
	var customers []model.Customer
	result := suite.db.Connection.Find(&customers)
	suite.NoError(result.Error)
	suite.Empty(customers)

	mockCognito := new(MockCognitoRemoteDataSource)
	repo := repository.NewCustomerRepository(suite.db, mockCognito)

	newCustomer := entity.Customer{
		Name:  "Teste",
		CPF:   "12312312312",
		Email: "teste@teste.com",
	}

	newCustomerModel := &model.Customer{
		Name:  "Teste",
		CPF:   "12312312312",
		Email: "teste@teste.com",
	}

	mockCognito.On("SignUp", newCustomerModel).Return(&responses.NetworkError{
		Code: 419,
	})

	newId, err := repo.CreateCustomer(suite.ctx, newCustomer)

	suite.Error(err)
	suite.Equal(uint(0), newId)
}

func (suite *RepositoryTestSuite) TestUpdateCustomerWithSuccess() {
	// ensure that the postgres database is empty
	var customers []model.Customer
	result := suite.db.Connection.Find(&customers)
	suite.NoError(result.Error)
	suite.Empty(customers)

	mockCognito := new(MockCognitoRemoteDataSource)
	repo := repository.NewCustomerRepository(suite.db, mockCognito)

	newCustomer := entity.Customer{
		Name:  "Teste",
		CPF:   "12312312312",
		Email: "teste@teste.com",
	}

	newCustomerModel := &model.Customer{
		Name:  "Teste",
		CPF:   "12312312312",
		Email: "teste@teste.com",
	}

	mockCognito.On("SignUp", newCustomerModel).Return(nil)

	newId, err := repo.CreateCustomer(suite.ctx, newCustomer)

	suite.NoError(err)
	suite.Equal(uint(1), newId)

	updateCustomerModel := entity.Customer{
		ID:    newId,
		Name:  "Teste 2",
		CPF:   "12312312312",
		Email: "teste@teste.com",
	}

	err = repo.UpdateCustomer(suite.ctx, updateCustomerModel)

	suite.NoError(err)
}

func (suite *RepositoryTestSuite) TestLoginWithSuccess() {
	mockCognito := new(MockCognitoRemoteDataSource)
	repo := repository.NewCustomerRepository(suite.db, mockCognito)

	mockCognito.On("Login", "123456").Return("TOKEN", nil)

	token, err := repo.Login(context.TODO(), "123456")

	suite.NoError(err)
	suite.Equal("TOKEN", token)
}

func (suite *RepositoryTestSuite) TestLoginWithCognitoError() {
	mockCognito := new(MockCognitoRemoteDataSource)
	repo := repository.NewCustomerRepository(suite.db, mockCognito)

	mockCognito.On("Login", "123456").Return("", &responses.NetworkError{
		Code: 401,
	})

	token, err := repo.Login(context.TODO(), "123456")

	suite.Error(err)
	suite.Empty(token)
}
