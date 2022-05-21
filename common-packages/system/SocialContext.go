package system

import (
	"database/sql"
	"fmt"

	"github.com/Masher828/MessengerBackend/common-packages/database"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

type SocialContextStruct struct {
	PostgresDB  *sql.DB
	MongoClient *mongo.Client
	Redis       *redis.Client
}

var SocialContext = SocialContextStruct{}

func PrepareSocialContext() error {
	psqldb, err := database.GetPostgresClient()
	if err != nil {
		fmt.Println("cannot connect to the postgres")
		return err
	}

	SocialContext.PostgresDB = psqldb

	redisdb, err := database.GetRedisClient()
	if err != nil {
		fmt.Println("cannot connect to the redis")
		return err
	}

	SocialContext.Redis = redisdb

	mongdb, err := database.GetMongoClient()
	if err != nil {
		fmt.Println("cannot connect to the mongo")
		return err
	}

	SocialContext.MongoClient = mongdb
	return nil
}
