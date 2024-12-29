package main

import (
	"log"

	"github.com/thalesmacedo1/covid-api/application/usecases"
	"github.com/thalesmacedo1/covid-api/config"
	"github.com/thalesmacedo1/covid-api/infrastructure/database/neo4j/repositories"
	"github.com/thalesmacedo1/covid-api/infrastructure/logger"
	"github.com/thalesmacedo1/covid-api/interfaces/api/handlers"
	"github.com/thalesmacedo1/covid-api/interfaces/routers"

	"github.com/thalesmacedo1/covid-api/infrastructure/database/neo4j"
)

func main() {
	// Carrega as configurações
	if err := config.LoadConfig(".env.example"); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Inicializa o logger
	logr := logger.NewLogrusLogger()

	// Inicializa o cliente Neo4j
	neo4jClient, err := neo4j.NewNeo4jClient(config.Settings.Neo4jURI, config.Settings.Neo4jUser, config.Settings.Neo4jPassword)
	if err != nil {
		logr.Fatalf("Failed to initialize Neo4j client: %v", err)
	}
	defer neo4jClient.Close()

	// Inicializa os repositórios
	countryRepo := repositories.NewNeo4jCountryRepository(neo4jClient.Driver)
	vaccineRepo := repositories.NewNeo4jVaccineRepository(neo4jClient.Driver)
	covidStatsRepo := repositories.NewNeo4jCovidStatsRepository(neo4jClient.Driver)
	vaccinationStatsRepo := repositories.NewNeo4jVaccinationStatsRepository(neo4jClient.Driver)

	// // Inicializa o cliente Redis
	// redisClient, err := redis.NewRedisClient(config.Settings.RedisHost, config.Settings.RedisPassword, config.Settings.RedisDB)
	// if err != nil {
	// 	logr.Fatalf("Failed to initialize Redis client: %v", err)
	// }
	// defer redisClient.Close()

	// Inicializa os use cases
	getCovidTotalsUC := usecases.NewGetCovidTotalsUseCase(covidStatsRepo)
	getVaccinatedPeopleUC := usecases.NewGetVaccinatedPeopleUseCase(vaccinationStatsRepo)
	getVaccinesUsedUC := usecases.NewGetVaccinesUsedUseCase(vaccineRepo)
	getCountryWithMostCasesUC := usecases.NewGetCountryWithMostCasesUseCase(covidStatsRepo, countryRepo)
	getMostUsedVaccineUC := usecases.NewGetMostUsedVaccineUseCase(vaccineRepo)

	// Inicializa os handlers
	covidHandler := handlers.NewCovidHandler(getCovidTotalsUC, getCountryWithMostCasesUC, logr)
	vaccinationHandler := handlers.NewVaccinationHandler(getVaccinatedPeopleUC, logr)
	vaccineHandler := handlers.NewVaccineHandler(getVaccinesUsedUC, getMostUsedVaccineUC, logr)

	// Configura o roteador usando Gin
	router := routers.Router(covidHandler, vaccinationHandler, vaccineHandler, logr)

	// Inicia o servidor HTTP
	logr.Info("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		logr.Fatalf("Server failed to start: %v", err)
	}
}
