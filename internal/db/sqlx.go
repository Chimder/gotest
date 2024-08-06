package db

import (
	"context"
	"log"

	"github.com/chimas/GoProject/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func DBConn(ctx context.Context) (*pgxpool.Pool, *sqlx.DB, error) {
	cfg, err := pgxpool.ParseConfig(config.LoadEnv().DB_URL)
	if err != nil {
		log.Fatalf("Unable to parse config: %v", err)
		return nil, nil, err
	}

	sqlcPool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
		return nil, nil, err
	}

	sqlxDB, err := sqlx.Connect("postgres", config.LoadEnv().DB_URL)
	if err != nil {
		log.Fatalf("Unable to connect to database using sqlx: %v", err)
		return nil, nil, err
	}

	return sqlcPool, sqlxDB, nil
}
