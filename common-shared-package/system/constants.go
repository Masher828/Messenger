package system

import "time"

//Controller Constants

const (
	//system constants

	AuthFailed                     = "AuthFAILED"
	AuthUserContext                = "UserContext"
	RequestStartTime               = "RequestStartTime"
	AccessTokenToUser              = "AccessTokenToUser:%s"
	DefaultHashSaltSize      int64 = 16
	DefaultAccessTokenExpiry       = time.Hour * 24 * 3
	ResetPasswordTokenKey          = "ResetPasswordToken:%s"
	ResetPasswordTokenExpiry       = time.Minute * 5
	MaxPasswordRetries             = 5
	//MongoDatabaseName constants

	MongoDatabaseName                = "messages_db"
	CollectionNameUser               = "users"
	CollectionNameConversation       = "conversations"
	CollectionNameConversationUnread = "conversation_unread"

	CollectionNameMessages = "messages"

	//ConversationType

	ConversationTypeOne2One       = "individual"
	ConversationTypeGroup         = "group"
	ConversationLimit       int64 = 20
	MessagesLimit           int64 = 50
)
