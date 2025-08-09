package model

type Passenger struct {
	PassengerID int      `json:"passengerId"`
	Survived    int      `json:"survived"`
	Pclass      int      `json:"pClass"`
	Name        string   `json:"name"`
	Sex         string   `json:"sex"`
	Age         *float64 `json:"age,omitempty"`
	SibSp       int      `json:"sibSp"`
	Parch       int      `json:"parch"`
	Ticket      string   `json:"ticket"`
	Fare        *float64 `json:"fare,omitempty"`
	Cabin       *string  `json:"cabin,omitempty"`
	Embarked    *string  `json:"embarked,omitempty"`
}

type FareHistogram struct {
	Percentiles []string `json:"percentiles"`
	Counts      []int    `json:"counts"`
}

// ErrorResponse represents a standard error message format.
type ErrorResponse struct {
	Message string `json:"message" example:"An error occurred"`
}
