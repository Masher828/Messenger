package system

import "errors"

var (
	ErrInternalServer     = errors.New("internal server error")
	ErrInvalidCredentials = errors.New("invalid email & password")
	ErrUserIsLockedOut    = errors.New("user is locked please try resetting the password")
	ErrEmailAlreadyExists = errors.New("email id already exists")
)

func getErrorMessageMap() map[error]bool {
	ErrorInResponseMap := map[error]bool{
		ErrInternalServer:     true,
		ErrInvalidCredentials: true,
		ErrUserIsLockedOut:    true,
		ErrEmailAlreadyExists: true,
	}
	return ErrorInResponseMap
}

func IsFunctionalError(err error) bool {
	errorMessageMap := getErrorMessageMap()
	return errorMessageMap[err]
}
