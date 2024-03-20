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
	userhabdler "user-managment-service/internal/http-server/handlers/user"
	"user-managment-service/internal/lib/logger"
	"user-managment-service/internal/lib/logger/sl"
	authservice "user-managment-service/internal/service/auth"
	userservice "user-managment-service/internal/service/user"
	"user-managment-service/internal/storage/storage/postgres"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New(cfg.Env)

	log.Debug("initializing server...", slog.String("addr", cfg.Address))

	// Storage
	storage, err := postgres.New(cfg.Storage)
	if err != nil {
		log.Error("failed to init storage", sl.Error(err))
		os.Exit(1)
	}
	log.Debug("storage initialized")

	// Cash
	cash, err := redis.New(cfg.Cash)
	if err != nil {
		log.Error("failed to init cash", sl.Error(err))
		os.Exit(1)
	}
	log.Debug("cash initialized")

	// Broker
	broker, err := rabbitmq.New(cfg.Broker)
	if err != nil {
		log.Error("failed to init message broker", sl.Error(err))
		os.Exit(1)
	}
	log.Debug("broker initialized")

	// Service layer
	authService := authservice.New(log, storage, cash, broker, cfg.Token)
	userService := userservice.New(log, storage)

	// Constroller layer
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	auth := authhandler.New(log, authService, cfg.Token)
	user := userhabdler.New(log, userService, cfg.Token)

	r.HandleFunc("/healthcheck", healthcheck.Register())
	// r.Route("/users", nil)
	r.Route("/auth", auth.Register())
	r.Route("/user", user.Register())

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
