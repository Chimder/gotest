package server

import (
	"net/http"

	"github.com/chimas/GoProject/internal/auth"
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
	// r.Use(httprate.LimitByIP(6, 10*time.Second))
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4000", "http://localhost:3000", "http://localhost", "https://magnetic-gabbi-chimas.koyeb.app/", "https://manka-next.vercel.app", "https://magnetic-gabbi-chimas.koyeb.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		MaxAge:           300,
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

	r.Post("/user/create", UserHandler.CreateOrCheckUser)
	r.Post("/user/cookie/delete", UserHandler.DeleteCookie)

	r.Group(func(r chi.Router) {
		// r.Use(httprate.Limit(3, time.Minute, httprate.WithKeyFuncs(
		// 	httprate.KeyByIP,
		// 	httprate.KeyByEndpoint,
		// )))
		r.Use(auth.AuthMiddleware)
		r.Route("/user", func(r chi.Router) {
			r.Get("/{email}", UserHandler.GetUser)
			r.Get("/session", UserHandler.GetSession)
			r.Get("/favorite/one", UserHandler.IsUserFavorite)
			r.Get("/favorite/list", UserHandler.UserFavList)
			r.Post("/toggle/favorite", UserHandler.ToggleFavorite)
			r.Delete("/delete", UserHandler.DeleteUser)
		})
	})

	return r
}
