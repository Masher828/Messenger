package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func GetMongoClient() (*mongo.Client, error) {
	options := options.Client()

	var timeOut = time.Second * 60

	options.ServerSelectionTimeout = &timeOut

	client, err := mongo.Connect(context.TODO(), options.ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, err
}
