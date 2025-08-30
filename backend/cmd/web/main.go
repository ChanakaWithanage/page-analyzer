package main

import (
	"context"
	"log"
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
		log.Printf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
