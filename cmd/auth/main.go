package main

import (
	"github.com/ElishaFlacon/fast-sobes-auth/internal/app"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/config"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/infrastructure"
)

func main() {
	cfg := config.NewConfig()
	log := infrastructure.NewLogger(cfg.LogBufSize)
	defer log.Stop()

	metrics := infrastructure.NewPrometheus(":7070", log)
	metrics.Run()
	defer log.Stop()

	a := app.NewApp(cfg, log)
	if err := a.Run(); err != nil {
		log.Fatal("Failed to stert app: %v", err)
	}

	// select {}
}
