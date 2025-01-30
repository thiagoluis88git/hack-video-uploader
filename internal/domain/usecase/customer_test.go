package usecase

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/entity"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/responses"
)

var (
	validateCPFUseCase = NewValidateCPFUseCase()
	saveCustomer       = entity.Customer{
		Name:  "Name",
		CPF:   "171.079.720-73",
		Email: "teste@teste.com",
	}

	mockedSaveCustomer = entity.Customer{
		Name:  "Name",
		CPF:   "17107972073",
		Email: "teste@teste.com",
	}

	customerById = entity.Customer{
		ID:    1,
		Name:  "Name",
		CPF:   "171.079.720-73",
		Email: "teste@teste.com",
	}

	customerByCPF = entity.Customer{
		ID:    1,
		Name:  "Name",
		CPF:   "070.732.860-83",
		Email: "teste@teste.com",
	}
)

func TestCustomerServices(t *testing.T) {
	t.Parallel()

	t.Run("got success when creating customer in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockCustomerRepository)
		sut := NewCreateCustomerUseCase(validateCPFUseCase, mockRepo)

		ctx := context.TODO()

		mockRepo.On("CreateCustomer", ctx, mockedSaveCustomer).Return(uint(1), nil)

		response, err := sut.Execute(ctx, saveCustomer)

		assert.NoError(t, err)
		assert.NotEmpty(t, response)

		assert.Equal(t, uint(1), response.Id)
	})

	t.Run("got error when creating customer in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockCustomerRepository)
		sut := NewCreateCustomerUseCase(validateCPFUseCase, mockRepo)

		ctx := context.TODO()

		mockRepo.On("CreateCustomer", ctx, mockedSaveCustomer).Return(uint(0), &responses.LocalError{
			Code:    responses.DATABASE_CONFLICT_ERROR,
			Message: "Conflict",
		})

		response, err := sut.Execute(ctx, saveCustomer)

		assert.Error(t, err)
		assert.Empty(t, response)

		var businessError *responses.BusinessResponse
		assert.Equal(t, true, errors.As(err, &businessError))
		assert.Equal(t, http.StatusConflict, businessError.StatusCode)
	})

	t.Run("got success when updating customer in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockCustomerRepository)
		sut := NewUpdateCustomerUseCase(validateCPFUseCase, mockRepo)

		ctx := context.TODO()

		mockRepo.On("UpdateCustomer", ctx, mockedSaveCustomer).Return(nil)

		err := sut.Execute(ctx, saveCustomer)

		assert.NoError(t, err)
	})

	t.Run("got error when updating customer in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockCustomerRepository)
		sut := NewUpdateCustomerUseCase(validateCPFUseCase, mockRepo)

		ctx := context.TODO()

		mockRepo.On("UpdateCustomer", ctx, mockedSaveCustomer).Return(&responses.NetworkError{
			Code:    404,
			Message: "Not Found",
		})

		err := sut.Execute(ctx, saveCustomer)

		assert.Error(t, err)

		var businessError *responses.BusinessResponse
		assert.Equal(t, true, errors.As(err, &businessError))
		assert.Equal(t, http.StatusNotFound, businessError.StatusCode)
	})

	t.Run("got success when login in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockCustomerRepository)
		sut := NewLoginCustomerUseCase(mockRepo)

		ctx := context.TODO()

		mockRepo.On("Login", ctx, "07073286083").Return("token", nil)

		response, err := sut.Execute(ctx, "07073286083")

		assert.NoError(t, err)
		assert.NotEmpty(t, response)
	})

	t.Run("got error when login in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockCustomerRepository)
		sut := NewLoginCustomerUseCase(mockRepo)

		ctx := context.TODO()

		mockRepo.On("Login", ctx, "07073286083").Return("", &responses.NetworkError{
			Code: 401,
		})

		response, err := sut.Execute(ctx, "07073286083")

		assert.Error(t, err)
		assert.Empty(t, response)
	})
}
