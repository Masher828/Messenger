package database

import (
	"context"

	"github.com/go-redis/redis"
)

func GetRedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong := client.Ping(context.TODO())
	if pong.Err() != nil {
		return nil, pong.Err()
	}
	return client, nil

}
