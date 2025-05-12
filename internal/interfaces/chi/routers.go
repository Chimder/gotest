package chi

import (
	"net/http"

	"github.com/chimas/GoProject/internal/auth"
	"github.com/chimas/GoProject/internal/interfaces/chi/handler"
	"github.com/chimas/GoProject/internal/repository"
	"github.com/chimas/GoProject/internal/service"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"

	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(db *pgxpool.Pool, rdb *redis.Client) http.Handler {

	r := chi.NewRouter()
	// r.Use(httprate.LimitByIP(6, 10*time.Second))
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4000", "http://localhost:3000", "http://localhost", "https://magnetic-gabbi-chimas.koyeb.app/", "https://manka-next.vercel.app", "https://magnetic-gabbi-chimas.koyeb.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.yaml")
	})

	r.Mount("/swagger/", httpSwagger.WrapHandler)

	r.Handle("/metrics", promhttp.Handler())
	repo := repository.NewRepository(db)

	userService := service.NewUserService(repo, rdb)
	mangaService := service.NewMangaService(repo, rdb)
	chapterService := service.NewChapterService(repo, rdb)

	userHandler := handler.NewUserHandler(userService)
	mangaHandler := handler.NewMangaHandler(mangaService)
	chapterHandler := handler.NewChapterHandler(chapterService)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	r.Route("/manga", func(r chi.Router) {
		r.Get("/", mangaHandler.Manga)
		r.Get("/many", mangaHandler.Mangas)
		r.Get("/popular", mangaHandler.Popular)
		r.Get("/chapter", chapterHandler.Chapter)
		r.Get("/filter", mangaHandler.Filter)
	})

	r.Post("/user/create", userHandler.CreateOrCheckUser)
	r.Get("/user/cookie/delete", userHandler.DeleteCookie)

	r.Group(func(r chi.Router) {
		// r.Use(httprate.Limit(3, time.Minute, httprate.WithKeyFuncs(
		// 	httprate.KeyByIP,
		// 	httprate.KeyByEndpoint,
		// )))
		r.Use(auth.AuthMiddleware)
		r.Route("/user", func(r chi.Router) {
			r.Get("/", userHandler.GetUser)
			r.Get("/session", userHandler.GetSession)
			r.Get("/favorite/one", userHandler.IsUserFavorite)
			r.Post("/toggle/favorite", userHandler.ToggleFavorite)
			r.Delete("/delete", userHandler.DeleteUser)
			r.Get("/favorite/list", userHandler.UserFavList)
		})
	})

	return r
}
