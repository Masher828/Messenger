package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func GetPostgresClient() (*sql.DB, error) {
	db, err := sql.Open("postgres", "port=5432 host=localhost user=postgres password=root@123 dbname=social_db sslmode=disable")
	if err != nil {
		return nil, err
	}
	return db, nil

}
