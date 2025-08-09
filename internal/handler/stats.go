package handler

import (
	"fmt"
	"github.com/dhope-nagesh/titanic-go-service/internal/model"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"gonum.org/v1/gonum/stat"
)

// GetFareHistogram godoc
// @Summary      Get fare price histogram
// @Description  Returns data for a bar chart of fare prices in percentiles
// @Tags         Statistics
// @Produce      json
// @Success      200  {object}  model.FareHistogram
// @Failure      500  {object}  model.ErrorResponse
// @Router       /stats/fare_histogram [get]
func (h *APIHandler) GetFareHistogram(c *gin.Context) {
	fares, err := h.Repo.GetFares()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to retrieve fares"})
		return
	}
	if len(fares) == 0 {
		c.JSON(http.StatusOK, model.FareHistogram{Percentiles: []string{}, Counts: []int{}})
		return
	}

	// The data must be sorted to calculate quantiles.
	sort.Float64s(fares)

	bins := 10
	percentilePoints := make([]float64, bins+1)
	labels := make([]string, bins)
	counts := make([]int, bins)

	// --- CORRECTED LOGIC ---
	// 1. Calculate the value at each 10th percentile (0th, 10th, 20th, etc.)
	for i := 0; i <= bins; i++ {
		quantile := float64(i) / float64(bins)
		// Call stat.Quantile for each individual quantile
		percentilePoints[i] = stat.Quantile(quantile, stat.Empirical, fares, nil)
	}

	// 2. Create labels for the histogram bins
	for i := 0; i < bins; i++ {
		labels[i] = fmt.Sprintf("%.2f - %.2f", percentilePoints[i], percentilePoints[i+1])
	}

	// 3. Count how many passengers fall into each bin
	fareIndex := 0
	for i := 0; i < bins; i++ {
		// The upper bound for this bin is the value of the next percentile point.
		upperBound := percentilePoints[i+1]
		count := 0
		// Go through the sorted fares and count how many fall into the current bin.
		for fareIndex < len(fares) && fares[fareIndex] <= upperBound {
			// Special case for the very first bin to include the 0th percentile value.
			if i == 0 {
				count++
				fareIndex++
				continue
			}
			// For all other bins, check if the fare is greater than the previous bin's upper bound.
			if fares[fareIndex] > percentilePoints[i] {
				count++
			}
			fareIndex++
		}
		counts[i] = count
	}

	c.JSON(http.StatusOK, model.FareHistogram{
		Percentiles: labels,
		Counts:      counts,
	})
}
