package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thalesmacedo1/covid-api/application/usecases"
	"github.com/thalesmacedo1/covid-api/infrastructure/logger"
)

type VaccinationHandler struct {
	GetVaccinatedPeopleUC usecases.GetVaccinatedPeopleUseCase
	Logger                logger.Logger
}

func NewVaccinationHandler(getVaccinatedPeopleUC usecases.GetVaccinatedPeopleUseCase, logger logger.Logger) *VaccinationHandler {
	return &VaccinationHandler{
		GetVaccinatedPeopleUC: getVaccinatedPeopleUC,
		Logger:                logger,
	}
}

// GetVaccinatedPeople godoc
// @Summary Retrieve vaccinated people data
// @Description Fetches the number of vaccinated people for a given country on a specific date
// @Tags Vaccination
// @Accept json
// @Produce json
// @Param countryCode path string true "ISO Country Code"
// @Param date path string true "Date in YYYY-MM-DD format" format(date)
// @Success 200 {object} usecases.GetVaccinatedPeopleOutput "Successful retrieval of vaccinated people data"
// @Failure 400 {object} map[string]string "Invalid date format"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/countries/{countryCode}/vaccinations/{date} [get]
func (h *VaccinationHandler) GetVaccinatedPeople(c *gin.Context) {
	countryCode := c.Param("countryCode")
	dateStr := c.Param("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		h.Logger.Warnf("Invalid date format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD."})
		return
	}

	input := usecases.GetVaccinatedPeopleInput{
		CountryCode: countryCode,
		Date:        date,
	}

	output, err := h.GetVaccinatedPeopleUC.Execute(c.Request.Context(), input)
	if err != nil {
		h.Logger.Errorf("Error executing GetVaccinatedPeopleUseCase: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve vaccinated people data."})
		return
	}

	c.JSON(http.StatusOK, output)
}
