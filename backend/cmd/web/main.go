package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/chanaka-withanage/page-analyzer/internal/analyzer"
	"github.com/chanaka-withanage/page-analyzer/internal/config"
	"github.com/chanaka-withanage/page-analyzer/internal/fetch"
	"github.com/chanaka-withanage/page-analyzer/internal/gateway"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.Load()

	f := fetch.New(cfg.FetchTimeout, cfg.MaxRedirects, cfg.MaxBytes)
	svc := analyzer.New(f)
	svc.SetDefaultTimeout(cfg.FetchTimeout)

	handler := gateway.NewMuxWithService(svc)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("server starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "err", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	slog.Info("shutdown signal received")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server shutdown failed", "err", err)
	} else {
		slog.Info("server shutdown completed")
	}
}
