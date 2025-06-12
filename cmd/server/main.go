package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"go-server/internal/config"
	"go-server/internal/eth"
	"go-server/internal/handlers"
	"go-server/internal/scheduler"
	"go-server/internal/ws"
)

func main() {
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ethc, err := eth.New(ctx, cfg)
	if err != nil {
		log.Fatalf("eth init: %v", err)
	}
	defer ethc.Close()

	hub := ws.NewHub()

	go scheduler.Run(ctx, ethc, time.Minute, hub)

	r := chi.NewRouter()
	handlers.Register(r, handlers.Deps{Cfg: cfg, Eth: ethc, Hub: hub})

	srv := &http.Server{Addr: ":" + cfg.HTTPPort, Handler: r}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-ctx.Done()
	shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(shutCtx)
}
