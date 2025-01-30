package usecase

import "github.com/klassmann/cpfcnpj"

type ValidateCPFUseCase struct{}

func NewValidateCPFUseCase() *ValidateCPFUseCase {
	return &ValidateCPFUseCase{}
}

func (usecase *ValidateCPFUseCase) Execute(cpf string) (string, bool) {
	cpf = cpfcnpj.Clean(cpf)
	return cpf, cpfcnpj.ValidateCPF(cpf)
}
