// Throwaway spike for issue #15: WebSocket echo server plus a page that
// heartbeats and counts drops, to test long-lived connections through a
// corporate proxy. Delete or move to docs after the validation ADR.
package main

import (
	"context"
	_ "embed"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/coder/websocket"

	"github.com/reh2ur/planning-poker/internal/config"
)

//go:embed index.html
var page []byte

func main() {
	cfg, err := config.FromOS()
	if err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}
	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           newHandler(slog.Default()),
		ReadHeaderTimeout: 5 * time.Second,
	}
	slog.Info("listening", "addr", cfg.HTTPAddr)
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func newHandler(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(page)
	})
	mux.HandleFunc("GET /ws", func(w http.ResponseWriter, r *http.Request) {
		echo(logger, w, r)
	})
	return mux
}

// echo never logs the remote address; only how long the connection lived and
// why it ended.
func echo(logger *slog.Logger, w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		logger.Warn("ws accept failed", "error", err)
		return
	}
	defer func() { _ = conn.CloseNow() }()
	conn.SetReadLimit(4096)

	start := time.Now()
	logger.Info("ws opened")
	err = pump(r.Context(), conn)
	logger.Info("ws closed",
		"duration", time.Since(start).Round(time.Second),
		"status", websocket.CloseStatus(err),
		"error", err)
}

func pump(ctx context.Context, conn *websocket.Conn) error {
	for {
		typ, r, err := conn.Reader(ctx)
		if err != nil {
			return err
		}
		w, err := conn.Writer(ctx, typ)
		if err != nil {
			return err
		}
		if _, err := io.Copy(w, r); err != nil {
			return err
		}
		if err := w.Close(); err != nil {
			return err
		}
	}
}
