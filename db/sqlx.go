package db

import (
	"log"

	"github.com/chimas/GoProject/config"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// type ctx context.Context
func DBConnection() (*sqlx.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	// ctx := context.Context
	db, err := sqlx.Connect("postgres", config.LoadEnv().DB_URL)

	if err != nil {
		log.Println("DBBEnv", config.LoadEnv().REDIS_URL)
		log.Fatal("Unable to connect to database:", err)
		return nil, err
	}

	return db, nil
}
