package routes

import (
	"github.com/Masher828/MessengerBackend/common-shared-package/system"
	"github.com/Masher828/MessengerBackend/messagesapp/controller"
	"github.com/gin-gonic/gin"
)

func PrepareRoutes(router *gin.Engine) {

	application := system.Application{}

	router.POST("/messenger/messages/send", application.Route(&controller.Controller{}, "SendMessage", false))
	//TODO Put a check if user is participant of that
	router.GET("/messenger/:conversationId/messages/list", application.Route(&controller.Controller{}, "GetMessagesForConversation", false))
	router.GET("/messenger/conversation/friend/:friendId", application.Route(&controller.Controller{}, "GetConversationWithFriend", false))
	router.GET("/messenger/conversation/:conversationId", application.Route(&controller.Controller{}, "GetConversationById", false))
	router.GET("/messenger/conversations", application.Route(&controller.Controller{}, "GetConversations", false))

}
