package handler

import (
	"github.com/dhope-nagesh/titanic-go-service/internal/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAllPassengers godoc
// @Summary      Get all passengers
// @Description  Returns a list of all passengers
// @Tags         Passengers
// @Produce      json
// @Success      200  {array}   model.Passenger
// @Failure      500  {object}  model.ErrorResponse
// @Router       /passengers [get]
func (h *APIHandler) GetAllPassengers(c *gin.Context) {
	passengers, err := h.Repo.GetAllPassengers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to retrieve passengers"})
		return
	}
	c.JSON(http.StatusOK, passengers)
}

// GetPassengerByID godoc
// @Summary      Get a passenger by ID
// @Description  Returns all data for a single passenger
// @Tags         Passengers
// @Produce      json
// @Param        id   path      int  true  "Passenger ID"
// @Success      200  {object}  model.Passenger
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /passengers/{id} [get]
func (h *APIHandler) GetPassengerByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Invalid passenger ID format"})
		return
	}

	passenger, err := h.Repo.GetPassengerByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, passenger)
}

// GetPassengerAttributes godoc
// @Summary      Get specific attributes for a passenger
// @Description  Returns only requested attributes for a passenger
// @Tags         Passengers
// @Produce      json
// @Param        id   path      int  true  "Passenger ID"
// @Param        attributes query []string true "List of attributes" collectionFormat(multi)
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /passengers/{id}/attributes [get]
func (h *APIHandler) GetPassengerAttributes(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Invalid passenger ID format"})
		return
	}

	attributes := c.QueryArray("attributes")
	if len(attributes) == 0 {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "You must provide at least one attribute."})
		return
	}

	passenger, err := h.Repo.GetPassengerByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Message: err.Error()})
		return
	}

	filteredData := filterPassengerAttributes(*passenger, attributes)
	c.JSON(http.StatusOK, filteredData)
}
