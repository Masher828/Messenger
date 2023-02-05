package system

import (
	"fmt"
	db "github.com/Masher828/MessengerBackend/common-shared-package/system/database"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

type TypeMessengerContext struct {
	MongoDB *mongo.Client
	Redis   *redis.Client
}

var MessengerContext = TypeMessengerContext{}

func PrepareMessengerContext() error {
	var err error = nil
	MessengerContext.Redis, err = db.ConnectRedis()
	if err != nil {
		fmt.Println("Error while connecting redis : ", err.Error())
		return err
	}

	MessengerContext.MongoDB, err = db.ConnectMongo()
	if err != nil {
		fmt.Println("Error while connecting mongo : ", err.Error())
		return err
	}

	return nil
}
