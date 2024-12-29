package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thalesmacedo1/covid-api/application/usecases"
	"github.com/thalesmacedo1/covid-api/infrastructure/logger"
)

type VaccineHandler struct {
	GetVaccinesUsedUC    usecases.GetVaccinesUsedUseCase
	GetMostUsedVaccineUC usecases.GetMostUsedVaccineUseCase
	Logger               logger.Logger
}

func NewVaccineHandler(getVaccinesUsedUC usecases.GetVaccinesUsedUseCase, getMostUsedVaccineUC usecases.GetMostUsedVaccineUseCase, logger logger.Logger) *VaccineHandler {
	return &VaccineHandler{
		GetVaccinesUsedUC:    getVaccinesUsedUC,
		GetMostUsedVaccineUC: getMostUsedVaccineUC,
		Logger:               logger,
	}
}

// GetVaccinesUsed godoc
// @Summary Retrieve vaccines used in a country
// @Description Fetches the list of vaccines used in a specific country
// @Tags Vaccine
// @Accept json
// @Produce json
// @Param countryCode path string true "ISO Country Code"
// @Success 200 {object} usecases.GetVaccinesUsedOutput "Successful retrieval of vaccines used"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/countries/{countryCode}/vaccines [get]
func (h *VaccineHandler) GetVaccinesUsed(c *gin.Context) {
	countryCode := c.Param("countryCode")

	input := usecases.GetVaccinesUsedInput{
		CountryCode: countryCode,
	}

	output, err := h.GetVaccinesUsedUC.Execute(c.Request.Context(), input)
	if err != nil {
		h.Logger.Errorf("Error executing GetVaccinesUsedUseCase: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve vaccines used."})
		return
	}

	c.JSON(http.StatusOK, output)
}

// GetMostUsedVaccine godoc
// @Summary Retrieve the most used vaccine in a region
// @Description Finds the most commonly used vaccine in a specified region
// @Tags Vaccine
// @Accept json
// @Produce json
// @Param regionName path string true "Region Name (e.g., South America, Europe, Asia)"
// @Success 200 {object} usecases.GetMostUsedVaccineOutput "Successful retrieval of most used vaccine"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/regions/{regionName}/vaccines/most-used [get]
func (h *VaccineHandler) GetMostUsedVaccine(c *gin.Context) {
	regionName := c.Param("regionName")

	input := usecases.GetMostUsedVaccineInput{
		RegionName: regionName,
	}

	output, err := h.GetMostUsedVaccineUC.Execute(c.Request.Context(), input)
	if err != nil {
		h.Logger.Errorf("Error executing GetMostUsedVaccineUseCase: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve most used vaccine."})
		return
	}

	c.JSON(http.StatusOK, output)
}
