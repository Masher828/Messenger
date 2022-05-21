package routes

import (
	"github.com/Masher828/MessengerBackend/common-packages/system"
	"github.com/Masher828/MessengerBackend/messagesapp/controllers"
	"github.com/zenazn/goji"
)

func PrepareRoutes(application *system.Application) {

	//conversation
	goji.Post("/messages/conversation", application.Route(&controllers.Controller{}, "CreateConversation"))
	goji.Get("/messages/conversation", application.Route(&controllers.Controller{}, "GetConversation"))
	goji.Get("/messages/converstaion/:conversationId", application.Route(&controllers.Controller{}, "GetConversationById"))

}
