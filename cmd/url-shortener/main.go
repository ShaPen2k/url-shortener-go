package main

import (
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/router"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/logger/sl/setup"
	"url-shortener/internal/lib/server"
	"url-shortener/internal/storage/sqlite"
)

func main() {
	// Init config
	cfg := config.ConfigLoad()

	// Init logger
	log := setup.SetupLogger(cfg.Env)
	log.Info("starting url-shortener", slog.String("env", cfg.Env))

	// Init storage
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	// Init router
	r := router.Setup(log, cfg.HTTPServer, storage)

	// Init HTTP server
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// Run server with graceful shutdown logic
	server.Run(log, srv, cfg.HTTPServer.Timeout)

	log.Info("server stopped")
}
