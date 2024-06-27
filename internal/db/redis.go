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
		log.Println("REdisEnv")
		panic(err)
	}

	opt.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	opt.PoolSize = 1000
	opt.MinIdleConns = 10

	return opt, nil
}
