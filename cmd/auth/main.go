package main

import (
	"github.com/ElishaFlacon/fast-sobes-auth/internal/app"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/config"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/infrastructure/logger"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/infrastructure/postgres"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/infrastructure/prometheus"
)

func main() {
	cfg := config.NewConfig()
	log := logger.NewLogger(cfg.Log.BufSize)
	defer log.Stop()

	metrics := prometheus.NewPrometheus(":7070", log)
	metrics.Run()
	defer log.Stop()

	db, err := postgres.NewPostgres((*postgres.Config)(cfg.Pg))
	log.Infof(
		"database config: host=%s port=%d user=%s pass=%s db=%s",
		cfg.Pg.Host,
		cfg.Pg.Port,
		cfg.Pg.User,
		cfg.Pg.DBName,
	)
	if err != nil {
		log.Fatal("failed to connect to database: %v", err)
	}

	a := app.NewApp(cfg, db, log)
	if err := a.Run(); err != nil {
		log.Fatal("failed to start app: %v", err)
	}

	// select {}
}
