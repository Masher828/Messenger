package routes

import (
	"github.com/Masher828/MessengerBackend/authapp/controller"
	"github.com/Masher828/MessengerBackend/common-shared-package/system"
	"github.com/gin-gonic/gin"
)

func PrepareRoutes(router *gin.Engine) {

	application := system.Application{}

	router.POST("/messenger/messages/send", application.Route(&controller.Controller{}, "SendMessage", false))
	router.POST("/messenger/:conversationId/messages/list", application.Route(&controller.Controller{}, "GetMessagesForConversation", false))

	router.POST("/messenger/conversations", application.Route(&controller.Controller{}, "GetOrCreateConversation", true))

}
