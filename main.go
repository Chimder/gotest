package main

import (
	"log"
	"net/http"

	"github.com/chimas/GoProject/config"
	"github.com/chimas/GoProject/db"
	_ "github.com/chimas/GoProject/docs"
	"github.com/chimas/GoProject/handler"
	"github.com/go-redis/redis/v9"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	mux := http.NewServeMux()
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4000"},
	})

	db, err := db.DBConnection()
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
		return
	}
	defer db.Close()

	opt, err := redis.ParseURL(config.LoadEnv().REDIS_URL)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(opt)

	handler := handler.NewMangaHandler(db, rdb)
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("GET /mangas", handler.Mangas)
	mux.HandleFunc("GET /manga/{name}", handler.Manga)
	mux.HandleFunc("GET /manga/{name}/{chapter}", handler.Chapter)
	mux.HandleFunc("GET /popular", handler.Popular)
	mux.HandleFunc("GET /filter", handler.Filter)
	// mux.HandleFunc("GET /search", handler.Search)

	log.Println("Listening...")
	http.ListenAndServe(":4000", c.Handler(mux))
}
