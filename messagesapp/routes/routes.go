package routes

import (
	"github.com/Masher828/MessengerBackend/common-packages/system"
	"github.com/Masher828/MessengerBackend/messagesapp/controllers"
	"github.com/zenazn/goji"
)

func PrepareRoutes(application *system.Application) {
	goji.Get("/create/user", application.Route(&controllers.Controller{}, "CreateUser"))
}
