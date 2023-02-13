package main

import (
	"fmt"
	"github.com/Masher828/MessengerBackend/common-shared-package/system"
	"github.com/Masher828/MessengerBackend/messagesapp/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	//logger.Core().With([]zap.Field{{Key: "test", String: "value"}})
	application := system.Application{}
	err := system.PrepareMessengerContext()
	if err != nil {
		fmt.Println("Error while preparing context : ", err.Error())
		return
	}

	app := gin.New()

	app.Use(application.PerformanceMeasure())

	routes.PrepareRoutes(app)

	app.Use(gin.Recovery())

	app.Use(application.ApplyAuth())

	port := ":8082"
	err = app.Run(port)
	if err != nil {
		fmt.Printf("Error while listening on port : %s\n", port)
		return
	}

}
