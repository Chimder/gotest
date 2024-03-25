package main

import (
	"log"
	"net/http"

	"github.com/chimas/GoProject/config"
	"github.com/chimas/GoProject/db"
	_ "github.com/chimas/GoProject/docs"
	"github.com/chimas/GoProject/handler"
	"github.com/chimas/GoProject/middleware"
	"github.com/go-redis/redis/v9"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	router := http.NewServeMux()
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

	handlerM := handler.NewMangaHandler(db, rdb)
	handlerU := handler.NewUserHandler(db, rdb)
	router.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	router.HandleFunc("GET /mangas", handlerM.Mangas)
	router.HandleFunc("GET /manga/{name}", handlerM.Manga)
	router.HandleFunc("GET /manga/{name}/{chapter}", handlerM.Chapter)
	router.HandleFunc("GET /popular", handlerM.Popular)
	router.HandleFunc("GET /filter", handlerM.Filter)
	router.HandleFunc("GET /user/{email}", handlerU.GetUser)
	// router.HandleFunc("POST /user/addfavorite",handler)
	// router.HandleFunc("DELETE /user",handler)

	server := http.Server{
		Addr:    ":4000",
		Handler: middleware.Logging(c.Handler(router)),
	}
	log.Println("Listening...")
	server.ListenAndServe()
}
