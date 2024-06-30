package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/chimas/GoProject/docs"
	"github.com/chimas/GoProject/internal/db"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	httpServer *http.Server
	pgdb       *sqlx.DB
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

	pgdb, err := db.DBConn()
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	opt, err := db.RedisCon()
	if err != nil {
		log.Fatal("Unable to connect to redis:", err)
	}
	rdb := redis.NewClient(opt)

	/////////////////////////////////////////////////////
	/////////////////////////////////////////////////////

	httpServer := &http.Server{
		Addr:         ":" + PORT,
		Handler:      NewRouter(pgdb, rdb),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	return &Server{
		httpServer: httpServer,
		pgdb:       pgdb,
		rdb:        rdb,
	}
}

func (s *Server) Close() {
	if s.pgdb != nil {
		s.pgdb.Close()
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
