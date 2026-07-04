package main

import (
	"embed"

	"github.com/ajanjairam/m3u-manager/cmd"
	"github.com/ajanjairam/m3u-manager/internal/database"
	"github.com/ajanjairam/m3u-manager/internal/repository"
	"github.com/ajanjairam/m3u-manager/internal/server"
	"github.com/gofiber/fiber/v3"
)

//go:embed assets/migrations/*.sql
var migrationsFS embed.FS

func main() {
	flags := cmd.RegisterFlags()
	db := database.NewDatabase()
	defer db.Close()
	db.Migrate(migrationsFS)

	repo := repository.New(db.DB)

	server.NewServer(fiber.Config{AppName: "M3U Manager"}).
		Cors(flags.Env).
		Register(db.DB, repo).
		Static(flags.Env).
		Start(flags.Host, flags.Port)
}
