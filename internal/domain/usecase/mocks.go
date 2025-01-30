package usecase

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/entity"
)

type MockOrderRepository struct {
	mock.Mock
}

type MockCustomerRepository struct {
	mock.Mock
}

type MockPaymentRepository struct {
	mock.Mock
}

type MockPaymentGatewayRepository struct {
	mock.Mock
}

type MockProductRepository struct {
	mock.Mock
}

type MockUserAdminRepository struct {
	mock.Mock
}

type MockQRCodePaymentRepository struct {
	mock.Mock
}

func (mock *MockCustomerRepository) CreateCustomer(ctx context.Context, customer entity.Customer) (uint, error) {
	args := mock.Called(ctx, customer)
	err := args.Error(1)

	if err != nil {
		return 0, err
	}

	return args.Get(0).(uint), nil
}

func (mock *MockCustomerRepository) Login(ctx context.Context, cpf string) (string, error) {
	args := mock.Called(ctx, cpf)
	err := args.Error(1)

	if err != nil {
		return "", err
	}

	return args.Get(0).(string), nil
}

func (mock *MockCustomerRepository) UpdateCustomer(ctx context.Context, customer entity.Customer) error {
	args := mock.Called(ctx, customer)
	err := args.Error(0)

	if err != nil {
		return err
	}

	return nil
}
