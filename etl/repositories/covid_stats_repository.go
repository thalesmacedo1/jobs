package repositories

import (
	"context"
	"log"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type CovidStatsRepository struct {
	driver neo4j.DriverWithContext
}

func NewCovidStatsRepository(driver neo4j.DriverWithContext) *CovidStatsRepository {
	return &CovidStatsRepository{driver: driver}
}

func (r *CovidStatsRepository) CreateCovidStats(ctx context.Context, countryCode string, date time.Time, newCases int, cumulativeCases int, newDeaths int, cumulativeDeaths int) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.Run(ctx,
		`MERGE (c:Country {code: $countryCode})
             MERGE (date:Date {date: $date})
             CREATE (s:CovidStats {
                 cumulativeCases: $cumulativeCases,
                 newCases: $newCases,
                 cumulativeDeaths: $cumulativeDeaths,
                 newDeaths: $newDeaths
             })
             MERGE (c)-[:REPORTED_ON]->(s)
             MERGE (s)-[:ON_DATE]->(date)`,
		map[string]interface{}{
			"countryCode":      countryCode,
			"date":             date,
			"cumulativeCases":  cumulativeCases,
			"newCases":         newCases,
			"cumulativeDeaths": cumulativeDeaths,
			"newDeaths":        newDeaths,
		})

	if err != nil {
		log.Printf("Error creating CovidStats for country %s on date %s: %v", countryCode, date.Format("2006-01-02"), err)
		return err
	}

	return err
}
