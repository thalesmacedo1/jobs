package routers

import (
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/thalesmacedo1/covid-api/infrastructure/logger"
	"github.com/thalesmacedo1/covid-api/interfaces/api/handlers"
	"github.com/thalesmacedo1/covid-api/interfaces/middleware"
)

func Router(covidHandler *handlers.CovidHandler, vaccinationHandler *handlers.VaccinationHandler, vaccineHandler *handlers.VaccineHandler, logger logger.Logger) *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware(logger))

	// use ginSwagger middleware to serve the API docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// COVID Endpoints

	// 1. Qual foi o total acumulado de casos e mortes de Covid-19 em um país específico em uma data determinada?
	router.GET("/api/v1/countries/:countryCode/covid/:date", covidHandler.GetTotals)

	// 4. Qual país registrou o maior número de casos acumulados até uma data específica?
	router.GET("/api/v1/countries/highest-cases?:date", covidHandler.GetCountryWithMostCases)

	// Vaccination Endpoints

	// 2. Quantas pessoas foram vacinadas com pelo menos uma dose em um determinado país em uma data específica?
	router.GET("/api/v1/countries/:countryCode/vaccinations/:date", vaccinationHandler.GetVaccinatedPeople)

	// Vaccine Endpoints

	// 3. Quais foram as vacinas usadas em um determinado país e em que data elas começaram a ser aplicadas?
	router.GET("/api/v1/countries/:countryCode/vaccines", vaccineHandler.GetVaccinesUsed)

	// 5. Qual foi a vacina mais utilizada em uma região específica?
	router.GET("/api/v1/regions/:regionName/vaccines/most-used", vaccineHandler.GetMostUsedVaccine)

	return router
}
