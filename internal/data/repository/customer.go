package repository

import (
	"context"

	"github.com/thiagoluis88git/hack-video-uploader/internal/data/model"
	"github.com/thiagoluis88git/hack-video-uploader/internal/data/remote"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/entity"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/repository"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/database"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"

	"gorm.io/gorm"
)

type CustomerRepositoryImpl struct {
	db            *database.Database
	cognitoRemote remote.CognitoRemoteDataSource
}

func NewCustomerRepository(db *database.Database, cognitoRemote remote.CognitoRemoteDataSource) repository.CustomerRepository {
	return &CustomerRepositoryImpl{
		db:            db,
		cognitoRemote: cognitoRemote,
	}
}

func (repository *CustomerRepositoryImpl) CreateCustomer(ctx context.Context, customer entity.Customer) (uint, error) {
	customerEntity := &model.Customer{
		Name:  customer.Name,
		CPF:   customer.CPF,
		Email: customer.Email,
	}

	err := repository.cognitoRemote.SignUp(customerEntity)

	if err != nil {
		return 0, responses.GetCognitoError(err)
	}

	err = repository.db.Connection.WithContext(ctx).Create(customerEntity).Error

	if err != nil {
		return 0, responses.GetDatabaseError(err)
	}

	return customerEntity.ID, nil
}

func (repository *CustomerRepositoryImpl) UpdateCustomer(ctx context.Context, customer entity.Customer) error {
	customerEntity := &model.Customer{
		Model: gorm.Model{ID: customer.ID},
		Name:  customer.Name,
		CPF:   customer.CPF,
		Email: customer.Email,
	}

	err := repository.db.Connection.WithContext(ctx).Save(&customerEntity).Error

	if err != nil {
		return responses.GetDatabaseError(err)
	}

	return nil
}

func (repository *CustomerRepositoryImpl) populateCustomer(customerEntity model.Customer) entity.Customer {
	return entity.Customer{
		ID:    customerEntity.ID,
		Name:  customerEntity.Name,
		CPF:   customerEntity.CPF,
		Email: customerEntity.Email,
	}
}

func (repository *CustomerRepositoryImpl) Login(ctx context.Context, cpf string) (string, error) {
	token, err := repository.cognitoRemote.Login(cpf)

	if err != nil {
		return "", responses.GetDatabaseError(err)
	}

	return token, nil
}
