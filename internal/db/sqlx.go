package db

import (
	"log"

	"github.com/chimas/GoProject/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func DBConn() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", config.LoadEnv().DB_URL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	return db, nil
}
