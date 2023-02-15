package main

import (
	"fmt"
	"github.com/Masher828/MessengerBackend/common-shared-package/conf"
	"github.com/Masher828/MessengerBackend/common-shared-package/system"
	"github.com/Masher828/MessengerBackend/messagesapp/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	application := system.Application{}

	err := conf.LoadConfigFile()
	if err != nil {
		fmt.Println("Error while loading config file : ", err.Error())
		return
	}

	err = system.PrepareMessengerContext()
	if err != nil {
		fmt.Println("Error while preparing context : ", err.Error())
		return
	}

	app := gin.New()

	app.Use(application.PerformanceMeasure())

	routes.PrepareRoutes(app)

	app.Use(gin.Recovery())

	app.Use(application.Cors())

	app.Use(application.ApplyAuth())

	port := ":8082"
	err = app.Run(port)
	if err != nil {
		fmt.Printf("Error while listening on port : %s\n", port)
		return
	}

}
