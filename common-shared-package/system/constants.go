package system

import "time"

//Controller Constants

const (
	AuthFailed                     = "AuthFAILED"
	AuthUserContext                = "UserContext"
	RequestStartTime               = "RequestStartTime"
	AccessTokenToUser              = "AccessTokenToUser:%s"
	DefaultHashSaltSize      int64 = 16
	DefaultAccessTokenExpiry       = time.Hour * 24 * 3

	MongoDatabaseName  = "messages_db"
	UserCollectionName = "users"
)
