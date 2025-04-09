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
