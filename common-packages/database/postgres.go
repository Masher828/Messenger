package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func GetPostgresClient() (*sql.DB, error) {
	host := viper.GetString("database.postgres.host")
	port := viper.GetString("database.postgres.port")
	user := viper.GetString("database.postgres.username")
	password := viper.GetString("database.postgres.password")
	dbname := viper.GetString("database.postgres.dbname")
	db, err := sql.Open("postgres", fmt.Sprintf("port=%s host=%s user=%s password=%s dbname=%s sslmode=disable", port, host, user, password, dbname))
	if err != nil {
		return nil, err
	}
	return db, nil

}
