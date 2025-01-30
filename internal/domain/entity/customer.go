package entity

type Customer struct {
	ID    uint   `json:"id"`
	Name  string `json:"name" validate:"required"`
	CPF   string `json:"cpf" validate:"required"`
	Email string `json:"email" validate:"required"`
}

type CustomerForm struct {
	CPF string `json:"cpf" validate:"required"`
}

type CustomerResponse struct {
	Id uint `json:"id"`
}
