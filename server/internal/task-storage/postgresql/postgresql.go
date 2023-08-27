package postgresql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/GrishaSkurikhin/Aviahackathon/internal/models"
	_ "github.com/lib/pq"
)

type TaskStorage struct {
	db *sql.DB
}

func New(host, port, user, password, dbname string) (*TaskStorage, error) {
	const op = "taskstorage.postgresql.New"

	info := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", info)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &TaskStorage{db: db}, nil
}

func (s *TaskStorage) GetTasks() ([]models.Task, error) {
	const op = "taskstorage.postgresql.GetTasks"

	stmt, err := s.db.Prepare("SELECT * FROM tasks WHERE status != 'done'")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task

		err := rows.Scan(&task.Id, &task.BusID, &task.FlightID, &task.TimeStart, &task.TimeEnd, &task.Status)
		if err != nil {
			return nil, fmt.Errorf("%s: scan statement: %w", op, err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *TaskStorage) GetBusTasks(driverID int) ([]models.Task, error) {
	const op = "taskstorage.postgresql.GetBusTasks"

	stmt, err := s.db.Prepare("SELECT * FROM tasks WHERE status != 'done' AND bus_id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(driverID)
	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task

		err := rows.Scan(&task.Id, &task.BusID, &task.FlightID, &task.TimeStart, &task.TimeEnd, &task.Status)
		if err != nil {
			return nil, fmt.Errorf("%s: scan statement: %w", op, err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *TaskStorage) ChangeTaskStatus(taskID int, newStatus string) error {
	const op = "taskstorage.postgresql.ChangeTaskStatus"

	stmt, err := s.db.Prepare("UPDATE tasks SET status = $1 WHERE id = $2")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(newStatus, taskID)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return nil
}

func (s *TaskStorage) ChangeTaskTime(taskID int, newTime time.Time) error {
	const op = "taskstorage.postgresql.ChangeTaskTime"

	stmt, err := s.db.Prepare("UPDATE tasks SET time_start = $1 WHERE id = $2")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(newTime.Format("2006-01-02 15:04:05"), taskID)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return nil
}

func (s *TaskStorage) ChangeTaskBus(taskID int, newBusID int) error {
	const op = "taskstorage.postgresql.ChangeTaskBus"

	stmt, err := s.db.Prepare("UPDATE tasks SET bus_id = $1 WHERE id = $2")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(newBusID, taskID)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return nil
}

func (s *TaskStorage) AddTasks(tasks []models.Task) error {
	const op = "taskstorage.postgresql.AddTasks"

	stmt, err := s.db.Prepare("INSERT INTO tasks (bus_id, flight_id, time_start, time_end, status) VALUES ($1, $2, $3, $4, %5)")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	for _, task := range tasks {
		_, err := stmt.Exec(task.BusID, task.FlightID, task.TimeStart, task.TimeEnd, task.Status)
		if err != nil {
			return fmt.Errorf("%s: execute statement: %w", op, err)
		}
	}

	return nil
}
