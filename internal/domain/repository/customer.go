package repository

import (
	"context"

	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/entity"
)

type CustomerRepository interface {
	CreateCustomer(ctx context.Context, customer entity.Customer) (uint, error)
	UpdateCustomer(ctx context.Context, customer entity.Customer) error
	Login(ctx context.Context, cpf string) (string, error)
}
