package usecase

import (
	"context"
	"net/http"

	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/entity"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/repository"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
)

type CreateCustomerUseCase interface {
	Execute(ctx context.Context, customer entity.Customer) (entity.CustomerResponse, error)
}

type CreateCustomerUseCaseImpl struct {
	validateCPFUseCase *ValidateCPFUseCase
	repository         repository.CustomerRepository
}

type UpdateCustomerUseCase interface {
	Execute(ctx context.Context, customer entity.Customer) error
}

type UpdateCustomerUseCaseImpl struct {
	validateCPFUseCase *ValidateCPFUseCase
	repository         repository.CustomerRepository
}

type LoginCustomerUseCase interface {
	Execute(ctx context.Context, cpf string) (entity.Token, error)
}

type LoginCustomerUseCaseImpl struct {
	repository repository.CustomerRepository
}

func NewUpdateCustomerUseCase(validateCPFUseCase *ValidateCPFUseCase, repository repository.CustomerRepository) UpdateCustomerUseCase {
	return &UpdateCustomerUseCaseImpl{
		validateCPFUseCase: validateCPFUseCase,
		repository:         repository,
	}
}

func NewCreateCustomerUseCase(validateCPFUseCase *ValidateCPFUseCase, repository repository.CustomerRepository) CreateCustomerUseCase {
	return &CreateCustomerUseCaseImpl{
		validateCPFUseCase: validateCPFUseCase,
		repository:         repository,
	}
}

func NewLoginCustomerUseCase(repository repository.CustomerRepository) LoginCustomerUseCase {
	return &LoginCustomerUseCaseImpl{
		repository: repository,
	}
}

func (service *CreateCustomerUseCaseImpl) Execute(ctx context.Context, customer entity.Customer) (entity.CustomerResponse, error) {
	cleanedCPF, validate := service.validateCPFUseCase.Execute(customer.CPF)

	if !validate {
		return entity.CustomerResponse{}, &responses.BusinessResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid CPF",
		}
	}

	customer.CPF = cleanedCPF
	customerId, err := service.repository.CreateCustomer(ctx, customer)

	if err != nil {
		return entity.CustomerResponse{}, responses.GetResponseError(err, "CustomerService")
	}

	return entity.CustomerResponse{
		Id: customerId,
	}, nil
}

func (service *UpdateCustomerUseCaseImpl) Execute(ctx context.Context, customer entity.Customer) error {
	cleanedCPF, validate := service.validateCPFUseCase.Execute(customer.CPF)

	if !validate {
		return &responses.BusinessResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid CPF",
		}
	}

	customer.CPF = cleanedCPF
	err := service.repository.UpdateCustomer(ctx, customer)

	if err != nil {
		return responses.GetResponseError(err, "CustomerService")
	}

	return nil
}

func (uc *LoginCustomerUseCaseImpl) Execute(ctx context.Context, cpf string) (entity.Token, error) {
	token, err := uc.repository.Login(ctx, cpf)

	if err != nil {
		return entity.Token{}, responses.GetResponseError(err, "CustomerService")
	}

	return entity.Token{
		AccessToken: token,
	}, nil
}
