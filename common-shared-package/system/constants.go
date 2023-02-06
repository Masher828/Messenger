package system

import "time"

//Controller Constants

const (
	//system constants

	AuthFailed                      = "AuthFAILED"
	AuthUserContext                 = "UserContext"
	RequestStartTime                = "RequestStartTime"
	AccessTokenToUser               = "AccessTokenToUser:%s"
	DefaultHashSaltSize       int64 = 16
	DefaultAccessTokenExpiry        = time.Hour * 24 * 3
	IncorrectPasswordCountKey       = "IncorrectPasswordCount:%s"
	MaxPasswordRetries              = 5
	// MongoDatabaseName constants

	MongoDatabaseName  = "messages_db"
	UserCollectionName = "users"
)
