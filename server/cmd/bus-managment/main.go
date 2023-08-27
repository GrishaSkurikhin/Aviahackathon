package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GrishaSkurikhin/Aviahackathon/internal/config"
	"github.com/GrishaSkurikhin/Aviahackathon/internal/lib/logger/sl"
	"github.com/GrishaSkurikhin/Aviahackathon/internal/lib/logger/slogpretty"
	"github.com/GrishaSkurikhin/Aviahackathon/internal/scheduler"
	"github.com/GrishaSkurikhin/Aviahackathon/internal/server"
	"golang.org/x/exp/slog"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"

	timeInterval = 30 // time between task generations in minutes
)

func main() {
	cfg := config.MustLoad()

	sched, err := scheduler.New(cfg, timeInterval*time.Minute)
	if err != nil {
		panic(err)
	}

	log := setupLogger(cfg.Env)

	log.Info(
		"starting bus-managment",
		slog.String("env", cfg.Env),
		slog.String("version", "2"),
	)
	log.Debug("debug messages are enabled")

	srv, err := server.New(cfg, log)
	if err != nil {
		log.Error("failed to create server", sl.Err(err))
	}
	log.Info("starting server", slog.String("address", cfg.HTTPServer.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Start(); err != nil {
			log.Error("failed to start server", sl.Err(err))
		}
	}()

	go func() {
		for {
			log.Info("Creating schedule")
			err := sched.Create()
			if err != nil {
				log.Error("failed to create schedule", sl.Err(err))
			} else {
				log.Info("Schedule created successfuly")
			}

			time.Sleep(timeInterval * time.Minute)
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Close(&ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
