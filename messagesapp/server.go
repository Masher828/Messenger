package main

import (
	"flag"
	"os"

	"github.com/Masher828/MessengerBackend/common-packages/system"
	"github.com/Masher828/MessengerBackend/messagesapp/routes"
	"github.com/zenazn/goji"
)

func main() {
	var application = &system.Application{}
	routes.PrepareRoutes(application)

	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	flag.Set("bind", "0.0.0.0:"+port)
	goji.Serve()
}
