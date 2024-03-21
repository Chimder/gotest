package db

import (
	"log"

	"github.com/chimas/GoProject/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func DBConnection(env config.EnvVars) *sqlx.DB {

	db := sqlx.MustConnect("postgress", env.DB_URL)

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully Connected")
	}

	return db
}
