package main

import (
	"fmt"
	"github.com/Masher828/MessengerBackend/authapp/routes"
	"github.com/Masher828/MessengerBackend/common-shared-package/conf"
	"github.com/Masher828/MessengerBackend/common-shared-package/system"
	"github.com/gin-gonic/gin"
)

func main() {

	err := conf.LoadConfigFile()
	if err != nil {
		fmt.Println("Error while loading conf file : ", err.Error())
		return
	}

	application := system.Application{}
	err = system.PrepareMessengerContext()
	if err != nil {
		fmt.Println("Error while preparing context : ", err.Error())
		return
	}

	app := gin.New()

	app.Use(application.PerformanceMeasure())

	app.Use(gin.Recovery())

	app.Use(application.Cors())

	app.Use(application.ApplyAuth())

	routes.PrepareRoutes(app)

	port := ":" + conf.MessengerConfig.Apps.Auth.Address
	err = app.Run(port)
	if err != nil {
		fmt.Printf("Error while listening on port : %s\n", port)
		return
	}

}
