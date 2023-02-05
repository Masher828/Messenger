package system

import "errors"

var (
	Err               = errors.New("hh")
	ErrInternalServer = errors.New("internal server error")
)

func getErrorMessageMap() map[error]bool {
	ErrorInResponseMap := map[error]bool{
		Err: false,
	}
	return ErrorInResponseMap
}

func IsFunctionalError(err error) bool {
	errorMessageMap := getErrorMessageMap()
	return errorMessageMap[err]
}
