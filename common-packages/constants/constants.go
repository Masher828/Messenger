package constants

var (

	//Conversation Type
	ConversationTypeGroup    = "GROUP"
	ConversationTypePersonal = "PERSONAL"

	//Mongo Database
	DatabaseSocialDB = "social_db"

	//Collection Name
	ConversationCollection     = "conversation"
	UserConversationCollection = "user_conversation"

	//Conversation defaults
	DefaultConversationOffset int64 = 0
	DefaultConversationLimit  int64 = 10
)
