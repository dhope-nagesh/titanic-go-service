package data

import "github.com/dhope-nagesh/titanic-go-service/internal/model"

// PassengerRepository defines the interface for data access.
type PassengerRepository interface {
	GetAllPassengers() ([]model.Passenger, error)
	GetPassengerByID(id int) (*model.Passenger, error)
	GetFares() ([]float64, error)
}
