package db

import (
	"context"
	"log"

	"github.com/chimas/GoProject/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func DBConn(ctx context.Context) (*pgx.Conn, *sqlx.DB, error) {
	sqlcDB, err := pgx.Connect(ctx, config.LoadEnv().DB_URL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		return nil, nil, err
	}

	sqlxDB, err := sqlx.Connect("postgres", config.LoadEnv().DB_URL)
	if err != nil {
		log.Fatalf("Unable to connect to database using sqlx: %v", err)
		return nil, nil, err
	}

	return sqlcDB, sqlxDB, nil
}
