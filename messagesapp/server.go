package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Masher828/MessengerBackend/common-packages/conf"
	"github.com/Masher828/MessengerBackend/common-packages/system"
	"github.com/Masher828/MessengerBackend/messagesapp/routes"
	"github.com/spf13/viper"
	"github.com/zenazn/goji"
)

func main() {
	err := conf.LoadConfigFile()
	if err != nil {
		fmt.Println(err.Error())
	}
	var application = system.Application{}
	routes.PrepareRoutes(&application)

	port := viper.GetString("apps.messagesapp.address")
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	flag.Set("bind", "0.0.0.0:"+port)
	goji.Serve()
}
