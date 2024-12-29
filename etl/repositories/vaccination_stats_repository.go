package repositories

import (
	"context"
	"log"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type VaccinationStatsRepository struct {
	driver neo4j.DriverWithContext
}

func NewVaccinationStatsRepository(driver neo4j.DriverWithContext) *VaccinationStatsRepository {
	return &VaccinationStatsRepository{driver: driver}
}

func (r *VaccinationStatsRepository) CreateVaccinationStats(
	ctx context.Context,
	countryCode string,
	totalVaccinations int,
	totalVaccinationsPer100 int,
	personsVaccinated1PlusDose int,
	personsVaccinated1PlusDosePer100 int,
	personsFullyVaccinated int,
	personsFullyVaccinatedPer100 int,
	personsBoosterAdditionalDose int,
	personsBoosterAdditionalDosePer100 int,
	dateUpdated time.Time,
) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.Run(ctx,
		`MATCH (c:Country {code: $countryCode})
         MERGE (date:Date {date: date($dateUpdated)})
         CREATE (v:VaccinationStats {
             totalVaccinations: $totalVaccinations,
             totalVaccinationsPer100: $totalVaccinationsPer100,
             personsVaccinated1PlusDose: $personsVaccinated1PlusDose,
             personsVaccinated1PlusDosePer100: $personsVaccinated1PlusDosePer100,
             personsFullyVaccinated: $personsFullyVaccinated,
             personsFullyVaccinatedPer100: $personsFullyVaccinatedPer100,
             personsBoosterAdditionalDose: $personsBoosterAdditionalDose,
             personsBoosterAdditionalDosePer100: $personsBoosterAdditionalDosePer100
         })
         MERGE (c)-[:VACCINATED_ON]->(v)
         MERGE (v)-[:ON_DATE]->(date)`,
		map[string]interface{}{
			"countryCode":                        countryCode,
			"dateUpdated":                        dateUpdated,
			"totalVaccinations":                  totalVaccinations,
			"totalVaccinationsPer100":            totalVaccinationsPer100,
			"personsVaccinated1PlusDose":         personsVaccinated1PlusDose,
			"personsVaccinated1PlusDosePer100":   personsVaccinated1PlusDosePer100,
			"personsFullyVaccinated":             personsFullyVaccinated,
			"personsFullyVaccinatedPer100":       personsFullyVaccinatedPer100,
			"personsBoosterAdditionalDose":       personsBoosterAdditionalDose,
			"personsBoosterAdditionalDosePer100": personsBoosterAdditionalDosePer100,
		})
	if err != nil {
		log.Printf("Error creating VaccinationStats for country %s on date %s: %v", countryCode, dateUpdated.Format("2006-01-02"), err)
	}
	return err
}
