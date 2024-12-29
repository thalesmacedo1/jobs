package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/thalesmacedo1/covid-api/domain/repositories"
	"github.com/thalesmacedo1/covid-api/domain/valueobjects"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var (
	dateFormat = "2006-01-02"
)

type Neo4jCovidStatsRepository struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jCovidStatsRepository(driver neo4j.DriverWithContext) repositories.CovidStatsRepository {
	return &Neo4jCovidStatsRepository{
		driver: driver,
	}
}

func (r *Neo4jCovidStatsRepository) GetTotalCasesAndDeaths(ctx context.Context, countryCode string, date time.Time) (*valueobjects.CovidStats, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	query := `
	MATCH (c:Country {code: $countryCode})-[:REPORTED_ON]->(cs:CovidStats)-[:ON_DATE]->(d:Date {date: $date})
	RETURN cs.cumulativeCases AS cumulativeCases, cs.newCases AS newCases, cs.cumulativeDeaths AS cumulativeDeaths, cs.newDeaths AS newDeaths
	`

	result, err := session.ExecuteRead(ctx, neo4j.ManagedTransactionWork(func(tx neo4j.ManagedTransaction) (interface{}, error) {
		rec, err := tx.Run(ctx, query, map[string]interface{}{
			"countryCode": countryCode,
			"date":        date.Format(dateFormat),
		})
		if err != nil {
			return nil, err
		}

		if rec.Next(ctx) {
			cumulativeCases, _ := rec.Record().Get("cumulativeCases")
			newCases, _ := rec.Record().Get("newCases")
			cumulativeDeaths, _ := rec.Record().Get("cumulativeDeaths")
			newDeaths, _ := rec.Record().Get("newDeaths")

			return valueobjects.NewCovidStats(
				int(cumulativeDeaths.(int64)),
				int(newDeaths.(int64)),
				int(cumulativeCases.(int64)),
				int(newCases.(int64)),
			), nil
		}

		if err = rec.Err(); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("no COVID-19 statistics found for country %s on date %s", countryCode, date.Format(dateFormat))
	}))

	if err != nil {
		return nil, err
	}

	return result.(*valueobjects.CovidStats), nil
}

func (r *Neo4jCovidStatsRepository) GetCountryWithMostCases(ctx context.Context, date time.Time) (string, int, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	query := `
	MATCH (c:Country)-[:REPORTED_ON]->(cs:CovidStats)-[:ON_DATE]->(d:Date {date: $date})
	RETURN c.code AS countryCode, cs.cumulativeCases AS cumulativeCases
	ORDER BY cs.cumulativeCases DESC
	LIMIT 1
	`

	result, err := session.ExecuteRead(ctx, neo4j.ManagedTransactionWork(func(tx neo4j.ManagedTransaction) (interface{}, error) {
		rec, err := tx.Run(ctx, query, map[string]interface{}{
			"date": date.Format(dateFormat),
		})
		if err != nil {
			return nil, err
		}

		if rec.Next(ctx) {
			countryCode, _ := rec.Record().Get("countryCode")
			cumulativeCases, _ := rec.Record().Get("cumulativeCases")
			return struct {
				Code  string
				Cases int
			}{
				Code:  countryCode.(string),
				Cases: int(cumulativeCases.(int64)),
			}, nil
		}

		if err = rec.Err(); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("no COVID-19 statistics found for date %s", date.Format(dateFormat))
	}))

	if err != nil {
		return "", 0, err
	}

	res := result.(struct {
		Code  string
		Cases int
	})

	return res.Code, res.Cases, nil
}
