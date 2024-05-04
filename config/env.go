package config

import (
	"os"
)

type EnvVars struct {
	REDIS_URL string
	DB_URL    string
}

func LoadEnv() EnvVars {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Println("Error loading .env file")
	// }

	redis_url := os.Getenv("REDIS_URL")
	db_url := os.Getenv("DB_URL")

	return EnvVars{
		REDIS_URL: redis_url,
		DB_URL:    db_url,
	}
}
