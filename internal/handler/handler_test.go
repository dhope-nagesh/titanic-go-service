package handler

import (
	"fmt"
	"testing"

	"github.com/dhope-nagesh/titanic-go-service/internal/model"
	"github.com/stretchr/testify/assert"
)

func ToPtr[T any](v T) *T {
	return &v
}

func TestFilterPassengerAttributes_AllAttributes(t *testing.T) {
	passenger := model.Passenger{
		PassengerID: 1,
		Name:        "John Doe",
		Age:         ToPtr(float64(30)),
		Sex:         "male",
		Pclass:      1,
		Survived:    1,
	}

	attrs := []string{"PassengerID", "Name", "Age", "Sex", "Pclass", "Survived"}
	result := filterPassengerAttributes(passenger, attrs)

	assert.Equalf(t, 1, result["passengerId"], fmt.Sprintf("Expected %v, got %v", 1, result))
	assert.Equal(t, "John Doe", result["name"])
	assert.Equal(t, ToPtr(float64(30)), result["age"])
	assert.Equal(t, "male", result["sex"])
	assert.Equal(t, 1, result["pClass"])
	assert.Equal(t, 1, result["survived"])
}

func TestFilterPassengerAttributes_SpecificAttributes(t *testing.T) {
	passenger := model.Passenger{
		PassengerID: 1,
		Name:        "John Doe",
		Age:         ToPtr(float64(30)),
		Sex:         "male",
		Pclass:      1,
		Survived:    1,
	}

	attrs := []string{"Name", "Age"}
	result := filterPassengerAttributes(passenger, attrs)

	assert.Equal(t, "John Doe", result["name"])
	assert.Equal(t, ToPtr(float64(30)), result["age"])
	assert.NotContains(t, result, "id")
	assert.NotContains(t, result, "sex")
}

func TestFilterPassengerAttributes_NoAttributes(t *testing.T) {
	passenger := model.Passenger{
		PassengerID: 1,
		Name:        "John Doe",
		Age:         ToPtr(float64(30)),
		Sex:         "male",
		Pclass:      1,
		Survived:    1,
	}

	attrs := []string{}
	result := filterPassengerAttributes(passenger, attrs)

	assert.Empty(t, result)
}

func TestFilterPassengerAttributes_InvalidAttributes(t *testing.T) {
	passenger := model.Passenger{
		PassengerID: 1,
		Name:        "John Doe",
		Age:         ToPtr(float64(30)),
		Sex:         "male",
		Pclass:      1,
		Survived:    0,
	}

	attrs := []string{"InvalidField"}
	result := filterPassengerAttributes(passenger, attrs)

	assert.Empty(t, result)
}
