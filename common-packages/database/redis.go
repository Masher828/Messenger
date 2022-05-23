package database

import (
	"context"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

func GetRedisClient() (*redis.Client, error) {

	host := viper.GetString("database.redisdb.host")
	port := viper.GetString("database.redisdb.port")
	password := viper.GetString("database.redisdb.password")
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})

	pong := client.Ping(context.TODO())
	if pong.Err() != nil {
		return nil, pong.Err()
	}

	return client, nil
}
