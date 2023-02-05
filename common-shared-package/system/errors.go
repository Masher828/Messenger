package system

import "errors"

var (
	ErrInternalServer     = errors.New("internal server error")
	ErrInvalidCredentials = errors.New("invalid email & password")
)

func getErrorMessageMap() map[error]bool {
	ErrorInResponseMap := map[error]bool{
		ErrInternalServer:     false,
		ErrInvalidCredentials: true,
	}
	return ErrorInResponseMap
}

func IsFunctionalError(err error) bool {
	errorMessageMap := getErrorMessageMap()
	return errorMessageMap[err]
}
