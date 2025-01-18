package identity

import "github.com/google/uuid"

type UUIDGenerator interface {
	New() string
}

type UUIDGeneratorImpl struct{}

func NewUUIDGenerator() UUIDGenerator {
	return &UUIDGeneratorImpl{}
}

func (gen *UUIDGeneratorImpl) New() string {
	return uuid.NewString()
}
