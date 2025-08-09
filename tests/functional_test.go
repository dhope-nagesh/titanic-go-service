package tests

import (
	"encoding/json"
	"github.com/dhope-nagesh/titanic-go-service/internal/data"
	"github.com/dhope-nagesh/titanic-go-service/internal/handler"
	"github.com/dhope-nagesh/titanic-go-service/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// setupFunctionalTestServer initializes a real repository and a real gin router
// for end-to-end testing.
func setupFunctionalTestServer(t *testing.T) *gin.Engine {
	// The test is run from the project root, so the path is relative to that.
	dbPath := "../data/titanic.db"
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatalf("Database file not found at %s. Please run 'make seed-sqlite' first.", dbPath)
	}

	// Use the real SQLite repository, not a mock.
	repo, err := data.NewSQLiteRepository(dbPath)
	assert.NoError(t, err)

	// Set up the router with the real handler.
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	apiHandler := handler.NewAPIHandler(repo)
	apiHandler.RegisterRoutes(router)

	return router
}

// TestFunctionalGetAllPassengers tests the happy path for retrieving all passengers.
func TestFunctionalGetAllPassengers(t *testing.T) {
	// Arrange
	router := setupFunctionalTestServer(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/passengers", nil)

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var passengers []model.Passenger
	err := json.Unmarshal(w.Body.Bytes(), &passengers)
	assert.NoError(t, err)
	assert.True(t, len(passengers) > 800, "Should return a full list of passengers")
	assert.Equal(t, "Braund, Mr. Owen Harris", passengers[0].Name, "The first passenger should match the data")
}

// TestFunctionalGetPassengerByID_Found tests retrieving a single, existing passenger.
func TestFunctionalGetPassengerByID_Found(t *testing.T) {
	// Arrange
	router := setupFunctionalTestServer(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/passengers/2", nil)

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var passenger model.Passenger
	err := json.Unmarshal(w.Body.Bytes(), &passenger)
	assert.NoError(t, err)
	assert.Equal(t, 2, passenger.PassengerID)
	assert.Equal(t, "Cumings, Mrs. John Bradley (Florence Briggs Thayer)", passenger.Name)
}

// TestFunctionalGetPassengerByID_NotFound tests the behavior for a non-existent passenger.
func TestFunctionalGetPassengerByID_NotFound(t *testing.T) {
	// Arrange
	router := setupFunctionalTestServer(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/passengers/9999", nil)

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestFunctionalGetPassengerAttributes tests the attribute filtering endpoint.
func TestFunctionalGetPassengerAttributes(t *testing.T) {
	// Arrange
	router := setupFunctionalTestServer(t)
	w := httptest.NewRecorder()
	// Build the request programmatically to ensure correct query parameter encoding.
	req := httptest.NewRequest("GET", "/api/v1/passengers/4/attributes", nil)
	q := req.URL.Query()
	q.Add("attributes", "Name")
	q.Add("attributes", "Sex")
	req.URL.RawQuery = q.Encode()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var attributes map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &attributes)
	assert.NoError(t, err)

	assert.Equal(t, "Futrelle, Mrs. Jacques Heath (Lily May Peel)", attributes["name"])
	assert.Equal(t, "female", attributes["sex"])
	assert.Nil(t, attributes["Age"], "Age should not be in the response")
	assert.Len(t, attributes, 2, "Response should only contain the 2 requested attributes")
}

// TestFunctionalGetFareHistogram tests the statistics endpoint.
func TestFunctionalGetFareHistogram(t *testing.T) {
	// Arrange
	router := setupFunctionalTestServer(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/stats/fare_histogram", nil)

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var histogram model.FareHistogram
	err := json.Unmarshal(w.Body.Bytes(), &histogram)
	assert.NoError(t, err)

	assert.NotEmpty(t, histogram.Counts)
	assert.NotEmpty(t, histogram.Percentiles)
	assert.Equal(t, 10, len(histogram.Counts), "Histogram should have 10 bins")
	assert.Equal(t, 10, len(histogram.Percentiles), "Histogram should have 10 labels")
}
