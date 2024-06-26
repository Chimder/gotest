package db

import (
	"crypto/tls"
	"log"

	"github.com/chimas/GoProject/internal/config"

	"github.com/redis/go-redis/v9"
)

func RedisCon() (*redis.Options, error) {

	opt, err := redis.ParseURL(config.LoadEnv().REDIS_URL)
	if err != nil {
		log.Println("REdisEnv", config.LoadEnv().REDIS_URL)
		panic(err)
	}

	opt.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	return opt, nil
}
