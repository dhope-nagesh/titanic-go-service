package data

import (
	"database/sql"
	"errors"
	"github.com/dhope-nagesh/titanic-go-service/internal/model"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &SQLiteRepository{db: db}, nil
}

func (r *SQLiteRepository) GetAllPassengers() ([]model.Passenger, error) {
	rows, err := r.db.Query("SELECT PassengerId, Survived, Pclass, Name, Sex, Age, SibSp, Parch, " +
		"Ticket, Fare, Cabin, Embarked FROM passengers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passengers []model.Passenger
	for rows.Next() {
		var p model.Passenger
		err := rows.Scan(&p.PassengerID, &p.Survived, &p.Pclass, &p.Name, &p.Sex, &p.Age, &p.SibSp, &p.Parch, &p.Ticket, &p.Fare, &p.Cabin, &p.Embarked)
		if err != nil {
			log.Printf("Error scanning passenger: %v", err)
			continue
		}
		passengers = append(passengers, p)
	}
	return passengers, nil
}

func (r *SQLiteRepository) GetPassengerByID(id int) (*model.Passenger, error) {
	row := r.db.QueryRow("SELECT PassengerId, Survived, Pclass, Name, Sex, Age, SibSp, Parch, Ticket, Fare, "+
		"Cabin, Embarked FROM passengers WHERE PassengerId = ?", id)

	var p model.Passenger
	err := row.Scan(&p.PassengerID, &p.Survived, &p.Pclass, &p.Name, &p.Sex, &p.Age, &p.SibSp, &p.Parch, &p.Ticket, &p.Fare, &p.Cabin, &p.Embarked)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("passenger not found")
		}
		return nil, err
	}
	return &p, nil
}

func (r *SQLiteRepository) GetFares() ([]float64, error) {
	rows, err := r.db.Query("SELECT Fare FROM passengers WHERE Fare IS NOT NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fares []float64
	for rows.Next() {
		var fare float64
		if err := rows.Scan(&fare); err != nil {
			continue
		}
		fares = append(fares, fare)
	}
	return fares, nil
}
