package scheduler

import (
	"fmt"
	"time"

	busstorage "github.com/GrishaSkurikhin/Aviahackathon/internal/bus-storage/postgresql"
	"github.com/GrishaSkurikhin/Aviahackathon/internal/config"
	flightstorage "github.com/GrishaSkurikhin/Aviahackathon/internal/flight-storage/postgresql"
	"github.com/GrishaSkurikhin/Aviahackathon/internal/models"
	distancegraph "github.com/GrishaSkurikhin/Aviahackathon/internal/models/distance-graph"
	taskstorage "github.com/GrishaSkurikhin/Aviahackathon/internal/task-storage/postgresql"
)

type FlightGetter interface {
	GetFlights(timeInterval time.Duration) ([]models.Flight, error)
}

type BusGetter interface {
	GetBuses() ([]models.Bus, error)
}

type TasksAdder interface {
	AddTasks(tasks []models.Task) error
}

type scheduler struct {
	flightGetter FlightGetter
	busGetter    BusGetter
	tasksAdder   TasksAdder
	timeInterval time.Duration
	distancegraph *distancegraph.Distancegraph
}

func New(cfg *config.Config, timeInterval time.Duration) (*scheduler, error) {
	const op = "lib.scheduler.New"

	flightGetter, err := flightstorage.New(cfg.FS.Host, cfg.FS.Port, cfg.FS.User, cfg.FS.Password, cfg.FS.DBname)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	busGetter, err := busstorage.New(cfg.BS.Host, cfg.BS.Port, cfg.BS.User, cfg.BS.Password, cfg.BS.DBname)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	tasksAdder, err := taskstorage.New(cfg.TS.Host, cfg.TS.Port, cfg.TS.User, cfg.TS.Password, cfg.TS.DBname)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	graph, err := distancegraph.New()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &scheduler{
		flightGetter: flightGetter,
		busGetter:    busGetter,
		tasksAdder:   tasksAdder,
		timeInterval: timeInterval,
		distancegraph: graph,
	}, nil
}

func (s *scheduler) Create() error {
	const op = "lib.scheduler.Create"

	flights, err := s.flightGetter.GetFlights(s.timeInterval)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	buses, err := s.busGetter.GetBuses()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	tasks := generateSchedule(flights, buses)

	err = s.tasksAdder.AddTasks(tasks)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}


func generateSchedule(flights []models.Flight, buses []models.Bus) []models.Task {
	const busCapacity = 30
	const busSpeed = 45 // km/h
	
	// TODO: finish algorithm

	// Идея алгоритма:

	// Цикл по всем активным автобусам: 
	//     Создаем переменные: список задач, время последней задачи автобуса, локация последней задачи автобуса
	//     Находим ближайший по времени рейс, на который успевает автобус (используем граф для подсчета расстояния и времени)
	//     Добавляем задачу автобусу, меняем его время и позицию.
	//     Из количества пассажиров рейса вычитаем количество пассажиров автобусов
	//	   Если пассажиров больше не осталось, больше не учитываем этот рейс
	//     Если нет задачи, удовлетворяющей автобусу, то переходим с следующему
	// Сложность алгоритма: O(n*m), где n - число автобусов, m - число рейсов
	return nil
}
