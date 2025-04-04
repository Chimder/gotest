package db

import (
	"crypto/tls"
	"time"

	"github.com/chimas/GoProject/internal/config"

	"github.com/redis/go-redis/v9"
)

func RedisConn() *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:         config.LoadEnv().REDIS_URL,
		DB:           0,
		PoolSize:     100,
		MinIdleConns: 10,
		MaxIdleConns: 50,
		PoolTimeout:  30 * time.Second,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})
	// opt, err := redis.ParseURL(config.LoadEnv().REDIS_URL)
	// if err != nil {
	// 	log.Println("REdisEnv")
	// 	panic(err)
	// }

	return client
}
