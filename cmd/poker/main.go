package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/reh2ur/planning-poker/internal/config"
	"github.com/reh2ur/planning-poker/internal/web"
)

func main() {
	cfg, err := config.FromOS()
	if err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}
	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           web.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	slog.Info("listening", "addr", cfg.HTTPAddr)
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
