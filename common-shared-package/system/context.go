package system

import (
	"fmt"

	db "github.com/Masher828/MessengerBackend/common-shared-package/system/database"
	internalMqtt "github.com/Masher828/MessengerBackend/common-shared-package/system/mqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

type TypeMessengerContext struct {
	MongoDB *mongo.Client
	Redis   *redis.Client
	Mqtt    mqtt.Client
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

	MessengerContext.Mqtt, err = internalMqtt.MqttConnect()
	if err != nil {
		fmt.Println("Error while connecting MQTT : ", err.Error())
		return err
	}

	fmt.Println(MessengerContext.Mqtt.Publish("/user/topic/339b2e5d-3a90-4792-9005-9d39fb0767ee", 0, false, "defef").Error())

	return nil
}
