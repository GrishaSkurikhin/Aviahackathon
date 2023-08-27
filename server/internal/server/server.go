package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/GrishaSkurikhin/Aviahackathon/internal/config"
	"github.com/GrishaSkurikhin/Aviahackathon/internal/server/handlers/tasks/change"
	"github.com/GrishaSkurikhin/Aviahackathon/internal/server/handlers/tasks/get"
	mwLogger "github.com/GrishaSkurikhin/Aviahackathon/internal/server/middleware/logger"

	taskstorage "github.com/GrishaSkurikhin/Aviahackathon/internal/task-storage/postgresql"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
)

type server struct {
	*http.Server
}

func New(cfg *config.Config, log *slog.Logger) (*server, error) {
	const op = "server.New"

	ts, err := taskstorage.New(cfg.TS.Host, cfg.TS.Port, cfg.TS.User, cfg.TS.Password, cfg.TS.DBname)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/get-tasks", func(r chi.Router) {
		r.Get("/", get.New(log, ts))
	})

	router.Route("/change-task", func(r chi.Router) {
		r.Post("/", change.New(log, ts))
	})

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	return &server{srv}, nil
}

func (srv *server) Start() error {
	const op = "server.Start"
	err := srv.ListenAndServe()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (srv *server) Close(ctx *context.Context) error {
	const op = "server.Close"
	err := srv.Shutdown(*ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
