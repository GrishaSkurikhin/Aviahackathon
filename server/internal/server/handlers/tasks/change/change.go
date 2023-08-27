package change

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	resp "github.com/GrishaSkurikhin/Aviahackathon/internal/lib/api/response"
	"github.com/GrishaSkurikhin/Aviahackathon/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

const (
	StatusWork     = "in work"
	StatusPause    = "on pause"
	StatusComplete = "complete"
	StatusQueue    = "queue"
)

type Request struct {
	TaskID    string `json:"taskID"`
	Parameter struct {
		Type  string `json:"type"` //status, time, busID
		Value string `json:"value"`
	} `json:"parameter"`
}

type TasksChanger interface {
	ChangeTaskStatus(int, string) error
	ChangeTaskTime(int, time.Time) error
	ChangeTaskBus(int, int) error
}

func New(log *slog.Logger, taskChanger TasksChanger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.tasks.change.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.JSON(w, r, resp.Error("empty request"))
			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		taskID, err := strconv.Atoi(req.TaskID)
		if err != nil {
			log.Error("wrong taskID format", sl.Err(err))
			render.JSON(w, r, resp.Error("wrong taskID format"))
			return
		}

		switch req.Parameter.Type {
		case "status":
			if req.Parameter.Value != StatusComplete && req.Parameter.Value != StatusPause &&
				req.Parameter.Value != StatusWork && req.Parameter.Value != StatusQueue {
				log.Error("wrong status")
				render.JSON(w, r, resp.Error("wrong status"))
				return
			}
			err := taskChanger.ChangeTaskStatus(taskID, req.Parameter.Value)
			if err != nil {
				log.Error("failed to change task", sl.Err(err))
				render.JSON(w, r, resp.Error("internal error"))
				return
			}

		case "time":
			layout := "2006-01-02 15:04:05"
			time, err := time.Parse(layout, req.Parameter.Value)
			if err != nil {
				log.Error("wrong time format", sl.Err(err))
				render.JSON(w, r, resp.Error("wrong time format"))
				return
			}
			err = taskChanger.ChangeTaskTime(taskID, time)
			if err != nil {
				log.Error("failed to change task", sl.Err(err))
				render.JSON(w, r, resp.Error("internal error"))
				return
			}

		case "busID":
			busID, err := strconv.Atoi(req.Parameter.Value)
			if err != nil {
				log.Error("wrong busID format", sl.Err(err))
				render.JSON(w, r, resp.Error("wrong busID format"))
				return
			}
			err = taskChanger.ChangeTaskBus(taskID, busID)
			if err != nil {
				log.Error("failed to change task", sl.Err(err))
				render.JSON(w, r, resp.Error("internal error"))
				return
			}

		default:
			log.Error("wrong parameter")
			render.JSON(w, r, resp.Error("wrong parameter"))
			return
		}

		log.Info("task changed correctly")
		render.JSON(w, r, resp.OK())
	}
}
