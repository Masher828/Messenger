package system

import (
	"database/sql"

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
		return err
	}
	SocialContext.PostgresDB = psqldb
	return nil
}
