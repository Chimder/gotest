package server

// import (
// 	"context"
// 	"log"
// 	"time"

// 	"github.com/chimas/GoProject/internal/db"
// 	"github.com/chimas/GoProject/internal/interfaces/fiber/handler"
// 	"github.com/chimas/GoProject/internal/repository"
// 	"github.com/chimas/GoProject/internal/service"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/gofiber/fiber/v2/middleware/compress"
// 	"github.com/gofiber/fiber/v2/middleware/logger"
// 	"github.com/gofiber/fiber/v2/middleware/recover"
// 	"github.com/jackc/pgx/v5/pgxpool"
// 	"github.com/redis/go-redis/v9"
// )

// type FiberServer struct {
// 	app *fiber.App
// 	db  *pgxpool.Pool
// 	rdb *redis.Client
// }

// func New() *FiberServer {
// 	app := fiber.New(fiber.Config{
// 		Prefork:       true,
// 		CaseSensitive: true,
// 		StrictRouting: true,
// 		ReadTimeout:   5 * time.Second,
// 		WriteTimeout:  10 * time.Second,
// 		IdleTimeout:   30 * time.Second,
// 	})

// 	app.Use(logger.New())
// 	app.Use(recover.New())
// 	app.Use(compress.New(compress.Config{
// 		Level: compress.LevelBestSpeed,
// 	}))

// 	pgdb, err := db.DBConn(context.Background())
// 	if err != nil {
// 		log.Fatal("Unable to connect to postgres:", err)
// 	}

// 	rdb := db.RedisConn()
// 	repo := repository.NewRepository(pgdb)

// 	userService := service.NewUserService(repo, rdb)
// 	mangaService := service.NewMangaService(repo, rdb)
// 	chapterService := service.NewChapterService(repo, rdb)

// 	userHandler := handler.NewUserHandler(userService)
// 	mangaHandler := handler.NewMangaHandler(mangaService)
// 	chapterHandler := handler.NewChapterHandler(chapterService)

// 	app.Get("/", func(c *fiber.Ctx) error {
// 		return c.SendString("Hello World")
// 	})
// 	app.Get("/api/list", func(c *fiber.Ctx) error {
// 		return c.SendString("I'm a GET request!")
// 	})
// 	app.Get("/yaml", func(c *fiber.Ctx) error {
// 		return c.SendFile("docs/swagger.yaml")
// 	})

// 	mangaR := app.Group("/manga")
// 	mangaR.Get("", mangaHandler.Manga)
// 	mangaR.Get("/many", mangaHandler.Mangas)
// 	mangaR.Get("/popular", mangaHandler.Popular)
// 	mangaR.Get("/chapter", chapterHandler.Chapter)

// 	userR := app.Group("/user")
// 	userR.Get("", userHandler.GetUser)
// 	userR.Get("/session", userHandler.GetSession)
// 	userR.Get("/favorite/one", userHandler.IsUserFavorite)
// 	userR.Post("/toggle/favorite", userHandler.ToggleFavorite)
// 	userR.Delete("/delete", userHandler.DeleteUser)
// 	userR.Get("/favorite/list", userHandler.UserFavList)

// 	return &FiberServer{
// 		app: app,
// 		db:  pgdb,
// 		rdb: rdb,
// 	}
// }

// func (s *FiberServer) Run() error {
// 	return s.app.Listen(":4000")
// }

// func (s *FiberServer) Shutdown(ctx context.Context) error {
// 	return s.app.ShutdownWithContext(ctx)
// }
