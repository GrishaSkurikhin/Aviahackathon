package models

import "time"

type Bus struct {
	Id      int
	Status  string
	Parking string
}

type Flight struct {
	Id          int
	Destination string
	Time        time.Time
	Status      string
	Passengers  int
}

type Task struct {
	Id        int       `json:"id"`
	BusID     int       `json:"busID"`
	FlightID  int       `json:"flightID"`
	TimeStart time.Time `json:"time start"`
	TimeEnd   time.Time `json:"time end"`
	Status    string    `json:"status"`
}
