package main

import (
	"log/slog"

	"github.com/deveusss/evergram-core/config"
	"github.com/deveusss/evergram-core/logging"
	"github.com/deveusss/evergram-identity/internal/app"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.Load[config.AppConfig]()
	log := setupLogger(cfg.Config.Env)
	app := app.NewApp(log, cfg.Config)
	app.MustRun()
}
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = logging.NewStdOutTextLogger(true, true)
	case envDev:
		log = logging.NewStdOutTextLogger(true, true)
	case envProd:
		log = logging.NewDefaultStdOutTextLogger()
	}

	return log
}
