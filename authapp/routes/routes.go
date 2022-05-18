package routes

import (
	"github.com/Masher828/MessengerBackend/authapp/controllers"
	"github.com/Masher828/MessengerBackend/common-packages/system"
	"github.com/zenazn/goji"
)

func PrepareRoutes(application *system.Application) {
	goji.Post("/auth/user/signup", application.Route(&controllers.Controller{}, "UserSignup"))
	goji.Post("/auth/user/signin", application.Route(&controllers.Controller{}, "UserSignin"))
	goji.Get("/auth/user/get", application.Route(&controllers.Controller{}, "GetUsers"))
}
