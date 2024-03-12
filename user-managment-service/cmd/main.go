package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"user-managment-service/internal/broker/rabbitmq"
	"user-managment-service/internal/cash/redis"
	"user-managment-service/internal/config"
	authhandler "user-managment-service/internal/http-server/handlers/auth"
	"user-managment-service/internal/http-server/handlers/healthcheck"
	"user-managment-service/internal/lib/logger"
	"user-managment-service/internal/lib/logger/sl"
	authservice "user-managment-service/internal/service/auth"
	"user-managment-service/internal/storage/storage/postgres"

	"github.com/go-chi/chi"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New(cfg.Env)

	log.Debug("initializing server...", slog.String("addr", cfg.Address))

	// Storage
	storage, err := postgres.New(cfg.Storage)
	if err != nil {
		log.Error("failed to init storage", sl.Error(err))
	}

	// Cash
	cash := redis.New(cfg.Cash)
	if cash == nil {
		log.Error("failed to init cash")
	}

	// Broker
	broker, err := rabbitmq.New(cfg.Broker)
	if err != nil {
		log.Error("failed to init message broker", sl.Error(err))
	}

	// Service layer
	authService := authservice.New(log, storage, cash, broker)

	// Constroller layer
	r := chi.NewRouter()

	auth := authhandler.New(log, authService, cfg.Token)

	r.HandleFunc("/healthcheck", healthcheck.Register())
	r.Route("/users", nil)
	r.Route("/auth", auth.Register())
	r.Route("/user", nil)

	// Server
	srv := http.Server{
		Handler:      r,
		Addr:         cfg.Address,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log.Debug("server initialized")
	log.Info("server is running...")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server", sl.Error(err))
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout*time.Second)
	defer cancel()

	srv.Shutdown(ctx)

	log.Info("server stopped")
}
