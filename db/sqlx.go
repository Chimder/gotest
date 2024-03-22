package db

import (
	"log"

	"github.com/chimas/GoProject/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func DBConnection() (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, "postgres", config.LoadEnv().DB_URL)

	if err != nil {
		log.Fatal("Unable to connect to database:", err)
		return nil, err
	}

	return db, nil
}
