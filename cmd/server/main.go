package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"room-decorator/internal/api"
	"room-decorator/internal/core"
	"room-decorator/internal/infra"
	"syscall"
	"time"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && a.Value.Kind() == slog.KindTime {
				a.Value = slog.TimeValue(a.Value.Time().UTC())
			}
			return a
		},
	})))

	repo := infra.NewInMemoryJobRepo()
	queue := infra.NewInMemoryQueue(10)

	go core.RunWorker(repo, queue)

	server := api.NewServer(repo, queue)

	httpServer := &http.Server{
		Addr:              ":8080",
		Handler:           server.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		slog.Info("listening", "addr", httpServer.Addr)
		serverErr <- httpServer.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	case <-ctx.Done():
		slog.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "err", err)
	}
	slog.Info("server stopped")
}
