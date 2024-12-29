package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/thalesmacedo1/covid-api/domain/repositories"
	"github.com/thalesmacedo1/covid-api/domain/valueobjects"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jVaccinationStatsRepository struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jVaccinationStatsRepository(driver neo4j.DriverWithContext) repositories.VaccinationStatsRepository {
	return &Neo4jVaccinationStatsRepository{
		driver: driver,
	}
}

func (r *Neo4jVaccinationStatsRepository) GetVaccinatedPeople(ctx context.Context, countryCode string, date time.Time) (*valueobjects.VaccinationStats, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	query := `
	MATCH (c:Country {code: $countryCode})-[:VACCINATED_ON]->(vs:VaccinationStats)-[:ON_DATE]->(d:Date {date: $date})
	RETURN vs.personsVaccinated1PlusDose AS personsVaccinated1PlusDose,
		   vs.personsVaccinated1PlusDosePer100 AS personsVaccinated1PlusDosePer100,
		   vs.totalVaccinations AS totalVaccinations,
		   vs.totalVaccinationsPer100 AS totalVaccinationsPer100
	`

	result, err := session.ExecuteRead(ctx, neo4j.ManagedTransactionWork(func(tx neo4j.ManagedTransaction) (interface{}, error) {
		rec, err := tx.Run(ctx, query, map[string]interface{}{
			"countryCode": countryCode,
			"date":        date.Format("2006-01-02"),
		})
		if err != nil {
			return nil, err
		}

		if rec.Next(ctx) {
			personsVaccinated1PlusDose, _ := rec.Record().Get("personsVaccinated1PlusDose")
			personsVaccinated1PlusDosePer100, _ := rec.Record().Get("personsVaccinated1PlusDosePer100")
			totalVaccinations, _ := rec.Record().Get("totalVaccinations")
			totalVaccinationsPer100, _ := rec.Record().Get("totalVaccinationsPer100")

			return valueobjects.NewVaccinationStats(
				int(personsVaccinated1PlusDosePer100.(int64)),
				int(personsVaccinated1PlusDose.(int64)),
				0, // Placeholder if not available
				0, // Placeholder if not available
				int(personsVaccinated1PlusDosePer100.(int64)),
				int(totalVaccinationsPer100.(int64)),
				int(personsVaccinated1PlusDose.(int64)),
				int(totalVaccinations.(int64)),
			), nil
		}

		if err = rec.Err(); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("no vaccination statistics found for country %s on date %s", countryCode, date.Format("2006-01-02"))
	}))

	if err != nil {
		return nil, err
	}

	return result.(*valueobjects.VaccinationStats), nil
}
