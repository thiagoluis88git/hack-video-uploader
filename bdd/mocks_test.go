package bdd_test

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/entity"
)

type MockGetTrackingsUseCase struct {
	mock.Mock
}

func (mock *MockGetTrackingsUseCase) Execute(ctx context.Context, cpf string) ([]entity.Tracking, error) {
	args := mock.Called(ctx, cpf)
	err := args.Error(1)

	if err != nil {
		return []entity.Tracking{}, err
	}

	return args.Get(0).([]entity.Tracking), nil
}
