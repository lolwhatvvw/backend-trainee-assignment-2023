package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lolwhatvvw/backend-trainee-assignment-2023/docs"
	"github.com/lolwhatvvw/backend-trainee-assignment-2023/internal/config"
	"github.com/lolwhatvvw/backend-trainee-assignment-2023/internal/handler"
	"github.com/lolwhatvvw/backend-trainee-assignment-2023/internal/router"
	"github.com/lolwhatvvw/backend-trainee-assignment-2023/internal/storage/postgres"
	httpswagger "github.com/swaggo/http-swagger/v2"
)

// @title Segment service API
// @version 1.0

const GracefulShutdownTimeout = 10 * time.Second

func main() {

	cfg := config.MustLoad()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	log.Info(fmt.Sprintf("starting %s on address %s and env=%s",
		cfg.Application.Name,
		fmt.Sprintf("%s", cfg.Server.Port),
		cfg.Env,
	))

	conn, err := postgres.NewConnection(cfg)

	userStorage, err := postgres.NewUserStorage(conn)
	if err != nil {
		log.Error(fmt.Sprintf("failed to connect to database %#v", cfg.Database), err)
	}

	segmentStorage, err := postgres.NewSegmentStorage(conn, userStorage)
	if err != nil {
		log.Error(fmt.Sprintf("failed to connect to database %#v", cfg.Database), err)
	}

	userController := handler.NewUserHandler(userStorage)
	segmentController := handler.NewSegmentHandler(segmentStorage)

	r := router.GetRouter(userController, segmentController)

	r.Get("/swagger/*", httpswagger.Handler(
		httpswagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", cfg.Server.Port)),
	))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("server shutting down")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("gracefully stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), GracefulShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server")
		os.Exit(1)
	}

	log.Info("server stopped")
}
