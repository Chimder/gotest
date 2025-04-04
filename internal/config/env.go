package config

import (
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type EnvVars struct {
	REDIS_URL            string `env:"REDIS_URL"`
	DB_URL               string `env:"DB_URL"`
	GOOGLE_CLIENT_ID     string `env:"GOOGLE_CLIENT_ID"`
	GOOGLE_CLIENT_SECRET string `env:"GOOGLE_CLIENT_SECRET"`
	IS_PROD              bool   `env:"IS_PROD"`
	COOKIE_SECRET        string `env:"COOKIE_SECRET"`
	PORT                 string `env:"PORT"`
}

func LoadEnv() *EnvVars {
	_ = godotenv.Load()
	cfg := EnvVars{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	return &cfg
}

// package config

// import (
// 	"log"
// 	"os"
// 	"strconv"

// 	"github.com/joho/godotenv"
// )

// type EnvVars struct {
// 	REDIS_URL            string
// 	DB_URL               string
// 	GOOGLE_CLIENT_ID     string
// 	GOOGLE_CLIENT_SECRET string
// 	IS_PROD              bool
// 	COOKIE_SECRET        string
// }

// func LoadEnv() EnvVars {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Println("Error loading .env file")
// 	}

// 	redis_url := os.Getenv("REDIS_URL")
// 	db_url := os.Getenv("DB_URL")
// 	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
// 	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

// 	cookieSecret := os.Getenv("COOKIE_SECRET")

// 	isProdStr := os.Getenv("IS_PROD")
// 	isProd, err := strconv.ParseBool(isProdStr)
// 	if err != nil {
// 		isProd = false
// 	}

// 	return EnvVars{
// 		REDIS_URL:            redis_url,
// 		DB_URL:               db_url,
// 		GOOGLE_CLIENT_ID:     googleClientId,
// 		GOOGLE_CLIENT_SECRET: googleClientSecret,
// 		IS_PROD:              isProd,
// 		COOKIE_SECRET:        cookieSecret,
// 	}
// }
