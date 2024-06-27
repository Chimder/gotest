package server

import (
	"net/http"

	"github.com/chimas/GoProject/internal/handler"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"

	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(pgdb *sqlx.DB, rdb *redis.Client) http.Handler {
	////////////////
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"http://localhost:4000", "http://localhost:3000", "http://localhost", "https://magnetic-gabbi-chimas.koyeb.app/", "https://manka-next.vercel.app", "https://magnetic-gabbi-chimas.koyeb.app"},
		AllowedOrigins: []string{"*"},
		// AllowCredentials: false,
		// MaxAge:           300,
	}))
	////////////////
	r.Get("/yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.yaml")
	})
	r.Mount("/swagger/", httpSwagger.WrapHandler)
	///////////////
	r.Handle("/metrics", promhttp.Handler())
	///////////////
	MangaHandler := handler.NewMangaHandler(pgdb, rdb)
	UserHandler := handler.NewUserHandler(pgdb, rdb)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	r.Route("/manga", func(r chi.Router) {
		r.Get("/", MangaHandler.Manga)
		r.Get("/many", MangaHandler.Mangas)
		r.Get("/chapter", MangaHandler.Chapter)
		r.Get("/popular", MangaHandler.Popular)
		r.Get("/filter", MangaHandler.Filter)
	})

	r.Route("/user", func(r chi.Router) {
		r.Get("/{email}", UserHandler.GetUser)
		r.Post("/create", UserHandler.CreateUserIfNotExists)
		r.Post("/favorite/{name}/{email}", UserHandler.ToggleFavorite)
		r.Get("/favorite/one", UserHandler.IsUserFavorite)
		r.Get("/favorite/list", UserHandler.UserFavList)
		r.Delete("/delete", UserHandler.DeleteUser)
	})

	return r
}
