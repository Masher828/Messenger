package database

import (
	"context"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func GetMongoClient() (*mongo.Client, error) {
	options := options.Client()

	var timeOut = time.Second * 60

	options.ServerSelectionTimeout = &timeOut

	host := viper.GetString("database.mongodb.host")
	port := viper.GetString("database.mongodb.port")

	client, err := mongo.Connect(context.TODO(), options.ApplyURI("mongodb://"+host+":"+port))
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, err
}
