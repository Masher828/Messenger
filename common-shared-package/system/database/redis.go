package db

import (
	"context"
	"github.com/Masher828/MessengerBackend/common-shared-package/conf"
	"github.com/go-redis/redis/v8"
)

func ConnectRedis() (*redis.Client, error) {

	port := conf.MessengerConfig.Database.Redis.Port
	password := conf.MessengerConfig.Database.Redis.Password
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1" + ":" + port,
		Password: password,
		DB:       0,
	})

	pong := client.Ping(context.TODO())
	if pong.Err() != nil {
		return nil, pong.Err()
	}

	return client, nil
}
