package services

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"path/filepath"

	"github.com/thalesmacedo1/covid-etl/config"
	"github.com/thalesmacedo1/covid-etl/repositories"
	"github.com/thalesmacedo1/covid-etl/utils"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type ETLService struct {
	driver               neo4j.DriverWithContext
	countryRepo          *repositories.CountryRepository
	vaccineRepo          *repositories.VaccineRepository
	covidStatsRepo       *repositories.CovidStatsRepository
	vaccinationStatsRepo *repositories.VaccinationStatsRepository
}

func NewETLService() (*ETLService, error) {
	driver, err := neo4j.NewDriverWithContext(
		config.Settings.Neo4jURI,
		neo4j.BasicAuth(config.Settings.Neo4jUser, config.Settings.Neo4jPassword, ""),
	)
	if err != nil {
		return nil, err
	}

	return &ETLService{
		driver:               driver,
		countryRepo:          repositories.NewCountryRepository(driver),
		vaccineRepo:          repositories.NewVaccineRepository(driver),
		covidStatsRepo:       repositories.NewCovidStatsRepository(driver),
		vaccinationStatsRepo: repositories.NewVaccinationStatsRepository(driver),
	}, nil
}

func (s *ETLService) Run(ctx context.Context) error {
	defer s.driver.Close(ctx)

	// // Create constraints
	// if err := s.createConstraints(ctx); err != nil {
	// 	return err
	// }

	// Process CSV files
	if err := s.processCountries(ctx, filepath.Join("data", "vaccination-data.csv")); err != nil {
		return err
	}
	if err := s.processVaccines(ctx, filepath.Join("data", "vaccination-metadata.csv")); err != nil {
		return err
	}
	if err := s.processCovidStats(ctx, filepath.Join("data", "WHO-COVID-19-global-data.csv")); err != nil {
		return err
	}

	return nil
}

// func (s *ETLService) createConstraints(ctx context.Context) error {
// 	session := s.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
// 	defer session.Close(ctx)

// 	constraints := []string{
// 		`CREATE CONSTRAINT IF NOT EXISTS FOR (n:Vaccine) REQUIRE (n.product) IS NODE KEY`,
// 		`CREATE CONSTRAINT IF NOT EXISTS FOR (n:Date) REQUIRE (n.date) IS NODE KEY`,
// 		`CREATE CONSTRAINT IF NOT EXISTS FOR (n:Region) REQUIRE (n.name) IS NODE KEY`,
// 		`CREATE CONSTRAINT IF NOT EXISTS FOR (n:Country) REQUIRE (n.code) IS NODE KEY`,
// 		`CREATE CONSTRAINT IF NOT EXISTS FOR (n:Vaccine) REQUIRE exists(n.company)`,
// 		`CREATE CONSTRAINT IF NOT EXISTS FOR (n:Vaccine) REQUIRE exists(n.vaccine)`,
// 		`CREATE CONSTRAINT IF NOT EXISTS FOR (n:Country) REQUIRE exists(n.name)`,
// 	}

// 	for _, constraint := range constraints {
// 		if _, err := session.Run(ctx, constraint, nil); err != nil {
// 			return fmt.Errorf("Error creating constraint %s: %v", constraint, err)
// 		}
// 	}
// 	log.Println("Constraints created successfully.")
// 	return nil
// }

func (s *ETLService) processCountries(ctx context.Context, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records[1:] {
		country := record[0]
		code := record[1]
		region := record[2]
		dateUpdated := record[4]

		// Create country entry
		if err := s.countryRepo.CreateCountry(ctx, country, code, region); err != nil {
			continue
		}

		// Create vaccination stats
		dateParsed := utils.ParseDate(dateUpdated)
		totalVaccinations := utils.ParseInt(record[5])
		totalVaccinationsPer100 := utils.ParseInt(record[7])
		personsVaccinated1PlusDose := utils.ParseInt(record[6])
		personsVaccinated1PlusDosePer100 := utils.ParseInt(record[8])
		personsLastDose := utils.ParseInt(record[9])
		personsLastDosePer100 := utils.ParseInt(record[10])
		personsBoosterAddDose := utils.ParseInt(record[14])
		personsBoosterAddDosePer100 := utils.ParseInt(record[15])

		if err := s.vaccinationStatsRepo.CreateVaccinationStats(
			ctx,
			code,
			totalVaccinations,
			totalVaccinationsPer100,
			personsVaccinated1PlusDose,
			personsVaccinated1PlusDosePer100,
			personsLastDose,
			personsLastDosePer100,
			personsBoosterAddDose,
			personsBoosterAddDosePer100,
			dateParsed); err != nil {
			continue
		}
	}
	return nil
}

func (s *ETLService) processVaccines(ctx context.Context, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';' // Set the correct delimiter

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records[1:] { // Skip header

		countryCode := record[0]
		product := record[1]
		vaccineName := record[2]
		company := record[3]
		authorizationDateStr := record[4]
		startDateStr := record[5]

		// Parse dates
		authorizationDate := utils.ParseDate(authorizationDateStr)
		startDate := utils.ParseDate(startDateStr)

		// Create vaccine
		err := s.vaccineRepo.CreateVaccine(ctx,
			product, vaccineName, company, authorizationDate, startDate)
		if err != nil {
			log.Printf("Error creating vaccine %s: %v", product, err)
			continue
		}

		// Create relationship between country and vaccine
		err = s.vaccineRepo.CreateCountryUsesVaccine(ctx, countryCode, product)
		if err != nil {
			log.Printf("Error creating relationship between country %s and vaccine %s: %v", countryCode, product, err)
			continue
		}
	}
	return nil
}

func (s *ETLService) processCovidStats(ctx context.Context, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';' // Set the correct delimiter

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records[1:] { // Skip header

		dateStr := record[0]
		countryCode := record[1]
		newCasesStr := record[4]
		cumulativeCasesStr := record[5]
		newDeathsStr := record[6]
		cumulativeDeathsStr := record[7]

		// Parse date and integers
		date := utils.ParseDate(dateStr)
		newCases := utils.ParseInt(newCasesStr)
		cumulativeCases := utils.ParseInt(cumulativeCasesStr)
		newDeaths := utils.ParseInt(newDeathsStr)
		cumulativeDeaths := utils.ParseInt(cumulativeDeathsStr)

		// Create COVID stats
		err := s.covidStatsRepo.CreateCovidStats(ctx,
			countryCode, date, newCases, cumulativeCases, newDeaths, cumulativeDeaths)
		if err != nil {
			log.Printf("Error creating Covid stats for country %s on date %s: %v", countryCode, dateStr, err)
			continue
		}
	}
	return nil
}
