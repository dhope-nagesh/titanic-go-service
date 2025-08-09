package data

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTempCSV(t *testing.T, content string) string {
	tmpFile, err := os.CreateTemp("", "testdata*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()
	return tmpFile.Name()
}

func TestNewCSVRepository(t *testing.T) {
	filePath := createTempCSV(t, "PassengerId,Survived,Pclass,Name,Sex,Age,SibSp,Parch,Ticket,Fare,Cabin,Embarked\n")
	defer os.Remove(filePath)

	repo, err := NewCSVRepository(filePath)
	assert.NoError(t, err)
	assert.NotNil(t, repo)
}

func TestCSVGetAllPassengers(t *testing.T) {
	filePath := createTempCSV(t, "PassengerId,Survived,Pclass,Name,Sex,Age,SibSp,Parch,Ticket,Fare,Cabin,Embarked\n1,1,1,John Doe,male,30,0,0,12345,100.0,C123,S\n")
	defer os.Remove(filePath)

	repo, err := NewCSVRepository(filePath)
	assert.NoError(t, err)

	passengers, err := repo.GetAllPassengers()
	assert.NoError(t, err)
	assert.Len(t, passengers, 1)
	assert.Equal(t, "John Doe", passengers[0].Name)
}

func TestCSVGetPassengerByID(t *testing.T) {
	filePath := createTempCSV(t, "PassengerId,Survived,Pclass,Name,Sex,Age,SibSp,Parch,Ticket,Fare,Cabin,Embarked\n1,1,1,John Doe,male,30,0,0,12345,100.0,C123,S\n")
	defer os.Remove(filePath)

	repo, err := NewCSVRepository(filePath)
	assert.NoError(t, err)

	passenger, err := repo.GetPassengerByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, passenger)
	assert.Equal(t, "John Doe", passenger.Name)
}

func TestCSVGetFares(t *testing.T) {
	filePath := createTempCSV(t, "PassengerId,Survived,Pclass,Name,Sex,Age,SibSp,Parch,Ticket,Fare,Cabin,Embarked\n1,1,1,John Doe,male,30,0,0,12345,100.0,C123,S\n2,0,3,Jane Doe,female,25,1,1,54321,200.0,,C\n")
	defer os.Remove(filePath)

	repo, err := NewCSVRepository(filePath)
	assert.NoError(t, err)

	fares, err := repo.GetFares()
	assert.NoError(t, err)
	assert.Len(t, fares, 2)
	assert.Equal(t, 100.0, fares[0])
	assert.Equal(t, 200.0, fares[1])
}
