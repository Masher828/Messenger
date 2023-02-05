package routes

import (
	"github.com/Masher828/MessengerBackend/authapp/controller"
	"github.com/Masher828/MessengerBackend/common-shared-package/system"
	"github.com/gin-gonic/gin"
)

func PrepareRoutes(router *gin.Engine) {

	application := system.Application{}
	//router.Group("/auth")
	router.GET("/signin", application.Route(&controller.Controller{}, "SignIn", true))
}
