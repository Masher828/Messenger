package db

import (
	"context"
	"fmt"
	"github.com/Masher828/MessengerBackend/common-shared-package/conf"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectMongo() (*mongo.Client, error) {

	host := conf.MessengerConfig.Database.Mongo.Host
	port := conf.MessengerConfig.Database.Mongo.Port
	prefix := conf.MessengerConfig.Database.Mongo.Prefix
	username := conf.MessengerConfig.Database.Mongo.Username
	password := conf.MessengerConfig.Database.Mongo.Password

	uri := prefix + "://" + username + ":" + password + "@" + host

	if len(port) != 0 {
		uri += ":" + port
	}

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)

	clientOptions := options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(serverAPIOptions)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	rp := readpref.Primary()
	err = client.Ping(context.TODO(), rp)
	if err != nil {
		fmt.Println(err)
	}

	return client, err
}
