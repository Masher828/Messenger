package routes

import (
	"github.com/Masher828/MessengerBackend/authapp/controller"
	"github.com/Masher828/MessengerBackend/common-shared-package/system"
	"github.com/gin-gonic/gin"
)

func PrepareRoutes(router *gin.Engine) {

	application := system.Application{}

	//authentication

	router.POST("/auth/signin", application.Route(&controller.Controller{}, "SignIn", true))
	router.POST("/auth/signup", application.Route(&controller.Controller{}, "SignUp", true))

	//profile

	router.POST("/auth/password/reset", application.Route(&controller.Controller{}, "ResetPassword", true))
	router.POST("/auth/password/update", application.Route(&controller.Controller{}, "UpdatePassword", true))
	router.PUT("/auth/profile", application.Route(&controller.Controller{}, "UpdateProfile", false))
	router.PUT("/auth/status/:status", application.Route(&controller.Controller{}, "UpdateStatus", false))
	router.GET("/auth/profile", application.Route(&controller.Controller{}, "GetProfile", false))
}
