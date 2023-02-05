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
	router.POST("/auth/password/otp/verify", application.Route(&controller.Controller{}, "VerifyOTP", true))
	router.POST("/auth/profile", application.Route(&controller.Controller{}, "UpdateProfile", false))
	router.GET("/auth/profile", application.Route(&controller.Controller{}, "GetProfile", false))
}
