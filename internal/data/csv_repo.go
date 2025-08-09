package data

import (
	"encoding/csv"
	"errors"
	"github.com/dhope-nagesh/titanic-go-service/internal/model"
	"io"
	"os"
	"strconv"
)

// CSVRepository holds the path to the CSV file.
type CSVRepository struct {
	filePath string
}

// NewCSVRepository creates a new instance of the CSV repository.
// It checks if the file exists before creating the repository.
func NewCSVRepository(filePath string) (*CSVRepository, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, errors.New("CSV file does not exist at the provided path: " + filePath)
	}
	return &CSVRepository{filePath: filePath}, nil
}

// read is a helper function to open, read, and parse the entire CSV file.
func (r *CSVRepository) read() ([]model.Passenger, error) {
	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Read() // Skip the header row

	var passengers []model.Passenger
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break // End of file
		}
		if err != nil {
			return nil, err
		}

		p, err := recordToPassenger(record)
		if err != nil {
			// In a real application, you might want to log this error
			// instead of stopping the entire process.
			continue
		}
		passengers = append(passengers, p)
	}
	return passengers, nil
}

// GetAllPassengers returns all passengers from the CSV file.
func (r *CSVRepository) GetAllPassengers() ([]model.Passenger, error) {
	return r.read()
}

// GetPassengerByID finds a single passenger by their ID in the CSV file.
func (r *CSVRepository) GetPassengerByID(id int) (*model.Passenger, error) {
	passengers, err := r.read()
	if err != nil {
		return nil, err
	}
	for _, p := range passengers {
		if p.PassengerID == id {
			return &p, nil
		}
	}
	return nil, errors.New("passenger not found")
}

// GetFares extracts all valid fare values from the CSV file.
func (r *CSVRepository) GetFares() ([]float64, error) {
	passengers, err := r.read()
	if err != nil {
		return nil, err
	}

	var fares []float64
	for _, p := range passengers {
		if p.Fare != nil {
			fares = append(fares, *p.Fare)
		}
	}
	return fares, nil
}

// recordToPassenger is a utility function to convert a single CSV record (a slice of strings)
// into a model.Passenger struct, handling type conversions and potential empty values.
func recordToPassenger(record []string) (model.Passenger, error) {
	var p model.Passenger
	var err error

	// PassengerId
	p.PassengerID, err = strconv.Atoi(record[0])
	if err != nil {
		return model.Passenger{}, errors.New("invalid PassengerId: " + record[0])
	}

	// Survived
	p.Survived, _ = strconv.Atoi(record[1])

	// Pclass
	p.Pclass, _ = strconv.Atoi(record[2])

	// Name & Sex
	p.Name = record[3]
	p.Sex = record[4]

	// Age (can be empty)
	if ageStr := record[5]; ageStr != "" {
		if age, err := strconv.ParseFloat(ageStr, 64); err == nil {
			p.Age = &age
		}
	}

	// SibSp & Parch
	p.SibSp, _ = strconv.Atoi(record[6])
	p.Parch, _ = strconv.Atoi(record[7])

	// Ticket
	p.Ticket = record[8]

	// Fare (can be empty)
	if fareStr := record[9]; fareStr != "" {
		if fare, err := strconv.ParseFloat(fareStr, 64); err == nil {
			p.Fare = &fare
		}
	}

	// Cabin (can be empty)
	if cabinStr := record[10]; cabinStr != "" {
		p.Cabin = &cabinStr
	}

	// Embarked (can be empty)
	if embarkedStr := record[11]; embarkedStr != "" {
		p.Embarked = &embarkedStr
	}

	return p, nil
}
