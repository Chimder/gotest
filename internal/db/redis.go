package db

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net/url"
	"time"

	"github.com/chimas/GoProject/internal/config"
	"github.com/redis/go-redis/v9"
)

func RedisConn() *redis.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg := config.LoadEnv()

	parsedURL, err := url.Parse(cfg.REDIS_URL)
	if err != nil {
		slog.Error("Invalid Redis URL", "err", err)
		return nil
	}

	password := ""
	if parsedURL.User != nil {
		password, _ = parsedURL.User.Password()
	}

	opts := &redis.Options{
		Addr:         parsedURL.Host,
		Password:     password,
		DB:           0,
		PoolSize:     100,
		MinIdleConns: 10,
		MaxIdleConns: 50,
		PoolTimeout:  30 * time.Second,
	}

	if cfg.IS_PROD {
		opts.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	client := redis.NewClient(opts)

	if _, err := client.Ping(ctx).Result(); err != nil {
		slog.Error("Failed to connect to Redis",
			"err", err,
			"host", parsedURL.Host,
			"is_prod", cfg.IS_PROD)
		return nil
	}

	slog.Info("Redis connected successfully", "host", parsedURL.Host, "is_prod", cfg.IS_PROD)

	return client
}
