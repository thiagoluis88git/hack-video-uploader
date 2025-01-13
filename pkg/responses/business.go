package responses

import (
	"errors"
	"fmt"
	"net/http"
)

type BusinessResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"msgError"`
}

func (br BusinessResponse) Error() string {
	return br.Message
}

/*
*

	Regras de validação de erros

	1) Verificação de erros de rede (BFF chamando serviços externos)
		Nesse caso o status code já vem do próprio serviço externo. Então só
		precisamos adaptar a mensagem de erro e usar o próprio status code para
		retornar pro usuário do BFF

	2) Verificação de erros de banco de dados (BFF chamando o banco do Microserviço)
		Nesse caso a mensagem de erro já vem do banco de dados (Ex: Duplicate keys). Então
		precisamos adaptar o status code e usar a própria mensagem de erro para
		retornar para o usuário do BFF

	3) Verificação de erros de outro Use Case (Use Case com dependência de outro Use Case)
		Nesse caso o status e a mensagem já estão prontas, é só repassar

	4) Default
		Caso não seja NetworkError ou LocalError, retornará um statuso code 500
		para o usuário

*
*/
func GetResponseError(err error, service string) error {
	var networkError *NetworkError
	var databaseError *LocalError
	var businessError *BusinessResponse

	statusCode := http.StatusInternalServerError
	message := "Unexpected internal error"

	if errors.As(err, &networkError) {
		statusCode = networkError.Code
		message = getBusinessMessageError(statusCode, service)
	} else if errors.As(err, &databaseError) {
		message = databaseError.Message
		statusCode = getBusinessStatusCode(*databaseError)
	} else if errors.As(err, &businessError) {
		statusCode = businessError.StatusCode
		message = businessError.Message
	}

	businessResponse := &BusinessResponse{
		StatusCode: statusCode,
		Message:    fmt.Sprintf("%v - %v", message, err.Error()),
	}

	return businessResponse
}

func getBusinessMessageError(statusCode int, service string) string {
	var message string

	switch statusCode {
	case http.StatusBadRequest:
		message = fmt.Sprintf("Bad request trying to execute %v", service)
	case http.StatusUnauthorized:
		message = fmt.Sprintf("Unauthorized error trying to execute %v", service)
	case http.StatusForbidden:
		message = fmt.Sprintf("Forbiden error trying to execute %v", service)
	case http.StatusNotFound:
		message = fmt.Sprintf("Not found trying to execute %v", service)
	case http.StatusConflict:
		message = fmt.Sprintf("Conflit with some data using the service %v", service)
	case http.StatusUnprocessableEntity:
		message = fmt.Sprintf("Logic error found in service %v", service)
	default:
		message = fmt.Sprintf("Unexpected internal error trying to execute service %v", service)
	}

	return message
}

func getBusinessStatusCode(localError LocalError) int {
	if localError.Code == DATABASE_CONFLICT_ERROR {
		return http.StatusConflict
	}

	if localError.Code == NOT_FOUND_ERROR {
		return http.StatusNotFound
	}

	if localError.Code == DATABASE_ERROR {
		return http.StatusServiceUnavailable
	}

	return http.StatusUnprocessableEntity
}
