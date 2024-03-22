package main

import (
	"log"
	"net/http"

	"github.com/chimas/GoProject/config"
	"github.com/chimas/GoProject/db"
	"github.com/chimas/GoProject/handler"
	"github.com/go-redis/redis/v9"
	_ "github.com/lib/pq"
)

func main() {
	mux := http.NewServeMux()
	// ctx := context.Background()

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
	// mux.HandleFunc("GET /", handler.Allmanga)
	mux.HandleFunc("GET /mangas", handler.Mangas)
	mux.HandleFunc("GET /manga/{name}", handler.Manga)
	mux.HandleFunc("GET /manga/{name}/{chapter}", handler.Chapter)
	mux.HandleFunc("GET /popular", handler.Popular)
	mux.HandleFunc("GET /search", handler.Search)
	mux.HandleFunc("GET /rating", handler.Rating)
	// mux.HandleFunc("GET /", handler.Allmanga)
	// mux.HandleFunc("GET /", handler.Allmanga)
	// mux.HandleFunc("GET /", handler.Allmanga)

	log.Println("Listening...")
	http.ListenAndServe(":5000", mux)
}
