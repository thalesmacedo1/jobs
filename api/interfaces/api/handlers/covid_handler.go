package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thalesmacedo1/covid-api/application/usecases"
	"github.com/thalesmacedo1/covid-api/infrastructure/logger"
)

type CovidHandler struct {
	GetCovidTotalsUC          usecases.GetCovidTotalsUseCase
	getCountryWithMostCasesUC usecases.GetCountryWithMostCasesUseCase
	Logger                    logger.Logger
}

func NewCovidHandler(getCovidTotalsUC usecases.GetCovidTotalsUseCase, getCountryWithMostCasesUC usecases.GetCountryWithMostCasesUseCase, logger logger.Logger) *CovidHandler {
	return &CovidHandler{
		GetCovidTotalsUC:          getCovidTotalsUC,
		getCountryWithMostCasesUC: getCountryWithMostCasesUC,
		Logger:                    logger,
	}
}

// GetTotals godoc
// @Summary Retrieve COVID-19 totals
// @Description Fetches COVID-19 total statistics for a specific country on a given date
// @Tags Covid
// @Accept json
// @Produce json
// @Param countryCode path string true "ISO Country Code"
// @Param date path string true "Date in YYYY-MM-DD format" format(date)
// @Success 200 {object} usecases.GetCovidTotalsOutput "Successful retrieval of COVID-19 totals"
// @Failure 400 {object} map[string]string "Invalid date format"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/countries/{countryCode}/covid/{date} [get]
func (h *CovidHandler) GetTotals(c *gin.Context) {
	countryCode := c.Param("countryCode")
	dateStr := c.Param("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		h.Logger.Warnf("Invalid date format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD."})
		return
	}

	input := usecases.GetCovidTotalsInput{
		CountryCode: countryCode,
		Date:        date,
	}

	output, err := h.GetCovidTotalsUC.Execute(c.Request.Context(), input)
	if err != nil {
		h.Logger.Errorf("Error executing GetCovidTotalsUseCase: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve COVID totals."})
		return
	}

	c.JSON(http.StatusOK, output)
}

// GetCountryWithMostCases godoc
// @Summary Retrieve country with most COVID-19 cases
// @Description Finds the country with the highest number of COVID-19 cases on a specific date
// @Tags Covid
// @Accept json
// @Produce json
// @Param date path string true "Date in YYYY-MM-DD format" format(date)
// @Success 200 {object} usecases.CountryWithMostCasesOutput "Successful retrieval of country with most cases"
// @Failure 400 {object} map[string]string "Invalid date format"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/countries/highest-cases?{date} [get]
func (h *CovidHandler) GetCountryWithMostCases(c *gin.Context) {
	dateStr := c.Param("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		h.Logger.Warnf("Invalid date format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD."})
		return
	}

	input := usecases.GetCountryWithMostCasesInput{
		Date: date,
	}

	country, err := h.getCountryWithMostCasesUC.Execute(c.Request.Context(), input)
	if err != nil {
		h.Logger.Errorf("Error getting country with most cases: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get country with most cases"})
		return
	}
	c.JSON(http.StatusOK, country)
}
