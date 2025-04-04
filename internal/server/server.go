package server

import (
	"context"
	"log"
	"net/http"
	"time"

	_ "github.com/chimas/GoProject/docs"
	"github.com/chimas/GoProject/internal/config"
	"github.com/chimas/GoProject/internal/db"
	"github.com/chimas/GoProject/internal/routers"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	httpServer *http.Server
	db         *pgxpool.Pool
	rdb        *redis.Client
}

func NewServer() *Server {
	ctx := context.Background()

	var PORT string
	if PORT = config.LoadEnv().PORT; PORT == "" {
		PORT = "4000"
	}
	///////////////
	recordMetrics()
	///////////////

	pgdb, err := db.DBConn(ctx)
	if err != nil {
		log.Fatal("Unable to connect to postgres:", err)
	}

	rdb := db.RedisConn()
	/////////////////////////////////////////////////////

	httpServer := &http.Server{
		Addr:         ":" + PORT,
		Handler:      routers.NewRouter(pgdb, rdb),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	return &Server{
		httpServer: httpServer,
		db:         pgdb,
		rdb:        rdb,
	}
}

func (s *Server) Close(ctx context.Context) {
	if s.db != nil {
		s.db.Close()
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
