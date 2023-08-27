package postgresql

import (
	"database/sql"
	"fmt"

	"github.com/GrishaSkurikhin/Aviahackathon/internal/models"
	_ "github.com/lib/pq"
)

type BusStorage struct {
	db *sql.DB
}

func New(host, port, user, password, dbname string) (*BusStorage, error) {
	const op = "busstorage.postgresql.New"

	info := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", info)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &BusStorage{db: db}, nil
}

func (s *BusStorage) GetBuses() ([]models.Bus, error) {
	const op = "busstorage.postgresql.GetBuses"

	stmt, err := s.db.Prepare("SELECT * FROM buses WHERE status = 'in work'")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	defer rows.Close()

	var buses []models.Bus
	for rows.Next() {
		var bus models.Bus

		err := rows.Scan(&bus.Id, &bus.Status, &bus.Parking)
		if err != nil {
			return nil, fmt.Errorf("%s: scan statement: %w", op, err)
		}
		buses = append(buses, bus)
	}

	return buses, nil
}
