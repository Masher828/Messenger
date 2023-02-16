package routes

import (
	"github.com/Masher828/MessengerBackend/common-shared-package/system"
	"github.com/Masher828/MessengerBackend/messagesapp/controller"
	"github.com/gin-gonic/gin"
)

func PrepareRoutes(router *gin.Engine) {

	application := system.Application{}

	router.POST("/messenger/messages/send", application.Route(&controller.Controller{}, "SendMessage", false))
	router.GET("/messenger/:conversationId/messages/list", application.Route(&controller.Controller{}, "GetMessagesForConversation", false))
	router.POST("/messenger/messages/:friendId", application.Route(&controller.Controller{}, "GetMessagesWithFriend", false))

	router.GET("/messenger/conversations", application.Route(&controller.Controller{}, "GetConversations", false))

}
