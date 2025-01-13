package responses

import (
	"net/http"
	"net/url"
	"strings"
)

type NetworkError struct {
	Code    int
	Message string
}

func (er NetworkError) Error() string {
	return er.Message
}

func GetCognitoError(err error) *NetworkError {
	code := http.StatusInternalServerError
	message := err.Error()

	if strings.Contains(err.Error(), "UsernameExistsException") {
		code = http.StatusConflict
	}

	return &NetworkError{
		Code:    code,
		Message: message,
	}
}

func GetNetworkError(urlError *url.Error) *NetworkError {
	code := http.StatusInternalServerError
	message := urlError.Unwrap().Error()

	if urlError.Timeout() {
		code = http.StatusGatewayTimeout
	}

	if urlError.Temporary() {
		code = http.StatusServiceUnavailable
	}

	return &NetworkError{
		Code:    code,
		Message: message,
	}
}

func IsNetworkResponseOk(response *http.Response, bodyMessage string) error {
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return nil
	}

	return &NetworkError{
		Code:    response.StatusCode,
		Message: bodyMessage,
	}
}
