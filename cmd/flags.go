package cmd

import (
	"flag"
	"log/slog"
	"os"

	"github.com/ajanjairam/m3u-manager/internal/models"
)

func RegisterFlags() models.FlagInput {
	var flags models.FlagInput

	flag.StringVar(&flags.Host, "host", "0.0.0.0", "Host to bind to")
	flag.Uint64Var(&flags.Port, "port", 5000, "Port to bind to")
	flag.StringVar(&flags.Env, "env", "production", "Environment to bind to")

	flag.Parse()

	if flags.Port > 65535 {
		slog.Error("Port must be between 0 and 65535")
		os.Exit(1)
	}
	return flags
}
