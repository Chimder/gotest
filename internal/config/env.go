package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvVars struct {
	REDIS_URL            string
	DB_URL               string
	GOOGLE_CLIENT_ID     string
	GOOGLE_CLIENT_SECRET string
	IS_PROD              bool
	COOKIE_SECRET        string
}

func LoadEnv() EnvVars {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	redis_url := os.Getenv("REDIS_URL")
	db_url := os.Getenv("DB_URL")
	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	cookieSecret := os.Getenv("COOKIE_SECRET")

	isProdStr := os.Getenv("IS_PROD")
	isProd, err := strconv.ParseBool(isProdStr)
	if err != nil {
		isProd = false
	}

	return EnvVars{
		REDIS_URL:            redis_url,
		DB_URL:               db_url,
		GOOGLE_CLIENT_ID:     googleClientId,
		GOOGLE_CLIENT_SECRET: googleClientSecret,
		IS_PROD:              isProd,
		COOKIE_SECRET:        cookieSecret,
	}
}
