package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/chimas/GoProject/docs"
	"github.com/chimas/GoProject/internal/db"
	"github.com/chimas/GoProject/internal/queries"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	httpServer *http.Server
	sqlc       *pgxpool.Pool
	sqlx       *sqlx.DB
	rdb        *redis.Client
}

func NewServer() *Server {
	var PORT string
	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "4000"
	}
	///////////////
	recordMetrics()
	///////////////

	ctx := context.Background()
	sqlc, sqlx, err := db.DBConn(ctx)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	sqlcQueries := queries.New(sqlc)

	opt, err := db.RedisCon()
	if err != nil {
		log.Fatal("Unable to connect to redis:", err)
	}
	rdb := redis.NewClient(opt)

	/////////////////////////////////////////////////////

	httpServer := &http.Server{
		Addr:         ":" + PORT,
		Handler:      NewRouter(sqlcQueries, sqlx, rdb),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	return &Server{
		httpServer: httpServer,
		sqlc:       sqlc,
		sqlx:       sqlx,
		rdb:        rdb,
	}
}

func (s *Server) Close(ctx context.Context) {
	if s.sqlc != nil {
		s.sqlc.Close()
	}
	if s.sqlx != nil {
		s.sqlx.Close()
	}
	if s.rdb != nil {
		s.rdb.Close()
	}
}

func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Addr() string {
	return s.httpServer.Addr
}
