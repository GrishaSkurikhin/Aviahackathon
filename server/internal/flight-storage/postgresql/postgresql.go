package postgresql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/GrishaSkurikhin/Aviahackathon/internal/models"
	_ "github.com/lib/pq"
)

type FlightStorage struct {
	db *sql.DB
}

func New(host, port, user, password, dbname string) (*FlightStorage, error) {
	const op = "flightstorage.postgresql.New"

	info := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", info)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &FlightStorage{db: db}, nil
}

func (s *FlightStorage) GetFlights(timeInterval time.Duration) ([]models.Flight, error) {
	const op = "flightstorage.postgresql.GetFlights"

	now := time.Now()
	endTime := now.Add(timeInterval)

	stmt, err := s.db.Prepare("SELECT * FROM flights WHERE time >= $1 AND time <= $2")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(now.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"))
	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	defer rows.Close()

	var flights []models.Flight
	for rows.Next() {
		var flight models.Flight

		err := rows.Scan(&flight.Id, &flight.Destination, &flight.Time, &flight.Status, &flight.Passengers)
		if err != nil {
			return nil, fmt.Errorf("%s: scan statement: %w", op, err)
		}
		flights = append(flights, flight)
	}

	return flights, nil
}
