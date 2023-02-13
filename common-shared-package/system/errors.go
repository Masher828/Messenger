package system

import "errors"

var (
	ErrInternalServer = errors.New("internal server error")

	//Auth

	ErrInvalidCredentials   = errors.New("invalid email & password")
	ErrUserIsLockedOut      = errors.New("user is locked please try resetting the password")
	ErrEmailAlreadyExists   = errors.New("email id already exists")
	ErrInvalidPasswordToken = errors.New("password token is not valid")
	ErrUnauthorizedAccess   = errors.New("got unauthorized user in protected routes")

	//Conversation

	ErrInvalidConversationType         = errors.New("invalid conversation type")
	ErrNotMemberOfConversation         = errors.New("user is not part of the conversation")
	ErrInvalidConversationId           = errors.New("invalid conversation id")
	ErrGroupConversationMinimumOneUser = errors.New("please add more users to create group conversation")
	ErrOne2OneConversationNoName       = errors.New("names are not assigned while creating one 2 one conversation")
)

func getErrorMessageMap() map[error]bool {
	ErrorInResponseMap := map[error]bool{
		ErrInternalServer:                  true,
		ErrInvalidCredentials:              true,
		ErrUserIsLockedOut:                 true,
		ErrEmailAlreadyExists:              true,
		ErrGroupConversationMinimumOneUser: true,
	}
	return ErrorInResponseMap
}

func IsFunctionalError(err error) bool {
	errorMessageMap := getErrorMessageMap()
	return errorMessageMap[err]
}
