package get

import (
	"net/http"
	"strconv"

	resp "github.com/GrishaSkurikhin/Aviahackathon/internal/lib/api/response"
	"github.com/GrishaSkurikhin/Aviahackathon/internal/lib/logger/sl"
	"github.com/GrishaSkurikhin/Aviahackathon/internal/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

type Response struct {
	resp.Response
	Tasks []models.Task `json:"tasks"`
}

type TasksGetter interface {
	GetTasks() ([]models.Task, error)
	GetBusTasks(int) ([]models.Task, error)
}

func New(log *slog.Logger, taskGetter TasksGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.tasks.get.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		busID := chi.URLParam(r, "busID")
		log.Debug("busID", busID)
		if busID == "" {
			log.Info("get all work tasks")

			tasks, err := taskGetter.GetTasks()
			if err != nil {
				log.Error("failed to get tasks", sl.Err(err))
				render.JSON(w, r, resp.Error("internal error"))
				return
			}
			render.JSON(w, r, ResponseOK(tasks))
		} else {
			log.Info("get work tasks for bus with id %s", busID)

			busID, err := strconv.Atoi(busID)
			if err != nil {
				log.Error("wrong parameter format", sl.Err(err))
				render.JSON(w, r, resp.Error("wrong parameter format"))
				return
			}

			tasks, err := taskGetter.GetBusTasks(busID)
			if err != nil {
				log.Error("failed to get tasks", sl.Err(err))
				render.JSON(w, r, resp.Error("internal error"))
				return
			}
			render.JSON(w, r, ResponseOK(tasks))
		}
		log.Info("tasks found and submitted")
	}
}

func ResponseOK(tasks []models.Task) Response {
	return Response{
		Response: resp.OK(),
		Tasks:    tasks,
	}
}
