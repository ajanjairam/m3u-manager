package server

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"

	"github.com/ajanjairam/m3u-manager/internal/repository"
	"github.com/ajanjairam/m3u-manager/internal/service"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/static"
)

type HTTPServer struct {
	app *fiber.App
}

func NewServer(cfg fiber.Config) *HTTPServer {
	app := fiber.New(cfg)
	return &HTTPServer{app: app}
}

func (s *HTTPServer) Cors(env string) *HTTPServer {
	if env == "dev" {
		s.app.Use(cors.New())
	}
	return s
}

func (s *HTTPServer) Register(db *sql.DB, repo *repository.Queries) *HTTPServer {
	v1Api := s.app.Group("/api/v1")
	v1Api.Get("/health-check", func(c fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"message": "The application is healthy!"})
	})

	ps := service.NewPlaylistService(db, repo)
	registerPlaylistRoutes(v1Api, ps)
	ts := service.NewChannelService(db, repo)
	registerChannelRoutes(v1Api, ts)
	fis := service.NewFilterService(db, repo)
	registerFilterRoutes(v1Api, fis)

	return s
}

func (s *HTTPServer) Static(env string, clientFS fs.FS) *HTTPServer {
	if env == "production" {
		clientSubFS, err := fs.Sub(clientFS, "client")
		if err != nil {
			log.Fatal(err)
		}

		s.app.Use(static.New("", static.Config{FS: clientSubFS}))
		s.app.Get("/*", func(c fiber.Ctx) error {
			data, err := fs.ReadFile(clientSubFS, "index.html")
			if err != nil {
				return c.SendStatus(fiber.StatusNotFound)
			}
			c.Type("html")
			return c.Send(data)
		})
	}
	return s
}

func (s *HTTPServer) Start(host string, port uint64) {
	log.Fatal(s.app.Listen(fmt.Sprintf("%s:%d", host, port)))
}
