package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chimas/GoProject/config"
	"github.com/chimas/GoProject/db"
	_ "github.com/chimas/GoProject/docs"
	"github.com/chimas/GoProject/handler"
	"github.com/chimas/GoProject/middleware"
	"github.com/go-redis/redis/v9"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "go_metrics",
		Subsystem: "prometheus",
		Name:      "processed_record_total",
		Help:      "process metrics count",
	})

	opsRequested = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "go_metrics",
		Subsystem: "prometheus",
		Name:      "processed_record_count",
		Help:      "request record count",
	})
)

func recordMetrics() {
	opsRequested.Inc()
	defer opsRequested.Dec()
	// loop
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

//		@title			Manka Api
//		@version		1.0
//		@description	Manga search
//	 @BasePath	/

func main() {

	recordMetrics()

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	router := http.NewServeMux()
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost", "http://localhost:4000", "http://localhost:3000", "https://magnetic-gabbi-chimas.koyeb.app/", "https://manka-next.vercel.app", "https://magnetic-gabbi-chimas.koyeb.app"},
	})

	db, err := db.DBConnection()
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
		return
	}
	defer db.Close()

	opt, err := redis.ParseURL(config.LoadEnv().REDIS_URL)
	if err != nil {
		log.Fatal("REdisEnv", err)
		return
	}
	/////////////
	opt.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	/////////////

	rdb := redis.NewClient(opt)

	handlerM := handler.NewMangaHandler(db, rdb)
	handlerU := handler.NewUserHandler(db, rdb)
	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("GET /yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.yaml")
	})
	router.HandleFunc("GET /app", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	router.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	router.HandleFunc("GET /mangas", handlerM.Mangas)
	router.HandleFunc("GET /manga", handlerM.Manga)
	router.HandleFunc("GET /manga/{name}/{chapter}", handlerM.Chapter)
	router.HandleFunc("GET /popular", handlerM.Popular)
	router.HandleFunc("GET /filter", handlerM.Filter)
	router.HandleFunc("GET /user/{email}", handlerU.GetUser)
	router.HandleFunc("POST /user/create", handlerU.CreateUserIfNotExists)
	router.HandleFunc("POST /user/favorite/{name}/{email}", handlerU.ToggleFavorite)
	router.HandleFunc("GET /user/favorite/one", handlerU.IsUserFavorite)
	router.HandleFunc("GET /user/favorite/list", handlerU.UserFavList)
	router.HandleFunc("DELETE /user/delete", handlerU.DeleteUser)

	var PORT string
	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "4000"
	}

	server := http.Server{
		Addr:    ":" + PORT,
		Handler: middleware.Logging(c.Handler(router)),
	}
	log.Println("Listening...")
	server.ListenAndServe()
}
