package data

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	return db, mock
}

func TestGetAllPassengers(t *testing.T) {
	// Mock database setup
	db, mock := setupMockDB(t)
	defer db.Close()

	mock.ExpectQuery("SELECT PassengerId, Survived, Pclass, Name, Sex, Age, SibSp, Parch, Ticket, Fare, Cabin, Embarked FROM passengers").
		WillReturnRows(sqlmock.NewRows([]string{"PassengerId", "Survived", "Pclass", "Name", "Sex", "Age", "SibSp", "Parch", "Ticket", "Fare", "Cabin", "Embarked"}).
			AddRow(1, 1, 1, "John Doe", "male", 30, 0, 0, "12345", 100.0, "C123", "S"))

	repo := &SQLiteRepository{db: db}
	passengers, err := repo.GetAllPassengers()

	assert.NoError(t, err)
	assert.Len(t, passengers, 1)
	assert.Equal(t, "John Doe", passengers[0].Name)
}

func TestGetPassengerByID(t *testing.T) {
	// Mock database setup
	db, mock := setupMockDB(t)
	defer db.Close()

	mock.ExpectQuery("SELECT PassengerId, Survived, Pclass, Name, Sex, Age, SibSp, Parch, Ticket, Fare, Cabin, Embarked FROM passengers WHERE PassengerId = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"PassengerId", "Survived", "Pclass", "Name", "Sex", "Age", "SibSp", "Parch", "Ticket", "Fare", "Cabin", "Embarked"}).
			AddRow(1, 1, 1, "John Doe", "male", 30, 0, 0, "12345", 100.0, "C123", "S"))

	repo := &SQLiteRepository{db: db}
	passenger, err := repo.GetPassengerByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, passenger)
	assert.Equal(t, "John Doe", passenger.Name)
}

func TestGetFares(t *testing.T) {
	// Mock database setup
	db, mock := setupMockDB(t)
	defer db.Close()

	mock.ExpectQuery("SELECT Fare FROM passengers WHERE Fare IS NOT NULL").
		WillReturnRows(sqlmock.NewRows([]string{"Fare"}).
			AddRow(100.0).
			AddRow(200.0))

	repo := &SQLiteRepository{db: db}
	fares, err := repo.GetFares()

	assert.NoError(t, err)
	assert.Len(t, fares, 2)
	assert.Equal(t, 100.0, fares[0])
	assert.Equal(t, 200.0, fares[1])
}
