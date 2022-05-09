package routes

import (
	"github.com/Masher828/MessengerBackend/authapp/controllers"
	"github.com/Masher828/MessengerBackend/common-packages/system"
	"github.com/zenazn/goji"
)

func PrepareRoutes(application *system.Application) {
	goji.Post("/auth/create/user", application.Route(&controllers.Controller{}, "CreateUser"))
}
