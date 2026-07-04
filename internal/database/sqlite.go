package database

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase() *Database {
	slog.Info("Initializing database")
	slog.Debug("Getting user home directory")
	var dbPath string
	userHome, err := os.UserHomeDir()
	if err != nil {
		dbPath = "app.db"
	} else {
		dir := filepath.Join(userHome, fmt.Sprintf(".%s", "m3u-manager"))
		info, err := os.Stat(dir)
		if err != nil || !info.IsDir() {
			if err = os.MkdirAll(dir, 0o755); err != nil {
				dbPath = "app.db"
			}
		}
		dbPath = filepath.Join(dir, "app.db")
	}
	slog.Info(fmt.Sprintf("Using database path: %s", dbPath))
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	return &Database{DB: db}
}

func (d *Database) Migrate(migrationsFS fs.FS) {
	subFS, err := fs.Sub(migrationsFS, "assets/migrations")
	if err != nil {
		log.Fatal(err)
	}
	provider, err := goose.NewProvider(goose.DialectSQLite3, d.DB, subFS)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	_, err = provider.Up(context.Background())
	if err != nil {
		slog.Error(err.Error())
	}
}

func (d *Database) Close() {
	err := d.DB.Close()
	if err != nil {
		log.Fatal(err)
	}
}
