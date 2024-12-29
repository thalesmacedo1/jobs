package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/thalesmacedo1/covid-api/domain/entities"
	"github.com/thalesmacedo1/covid-api/domain/repositories"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jCountryRepository struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jCountryRepository(driver neo4j.DriverWithContext) repositories.CountryRepository {
	return &Neo4jCountryRepository{
		driver: driver,
	}
}

func (r *Neo4jCountryRepository) GetCountryByCode(ctx context.Context, code string) (*entities.Country, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	query := `
	MATCH (c:Country {code: $code})
	RETURN c.name AS name, c.code AS code
	`

	result, err := session.ExecuteRead(ctx, neo4j.ManagedTransactionWork(func(tx neo4j.ManagedTransaction) (interface{}, error) {
		rec, err := tx.Run(ctx, query, map[string]interface{}{
			"code": strings.ToUpper(strings.TrimSpace(code)),
		})
		if err != nil {
			return nil, err
		}

		if rec.Next(ctx) {
			name, _ := rec.Record().Get("name")
			code, _ := rec.Record().Get("code")
			return entities.NewCountry(code.(string), name.(string))
		}

		if err = rec.Err(); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("country with code %s not found", code)
	}))

	if err != nil {
		return nil, err
	}

	return result.(*entities.Country), nil
}

func (r *Neo4jCountryRepository) CreateCountry(ctx context.Context, country *entities.Country) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close(ctx)

	query := `
	CREATE (c:Country {name: $name, code: $code})
	`

	_, err := session.ExecuteWrite(ctx, neo4j.ManagedTransactionWork(func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, query, map[string]interface{}{
			"name": country.Name,
			"code": strings.ToUpper(strings.TrimSpace(country.Code)),
		})
		return nil, err
	}))

	return err
}

func (r *Neo4jCountryRepository) AssociateCountryWithRegion(ctx context.Context, countryCode, regionName string) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close(ctx)

	query := `
	MATCH (c:Country {code: $countryCode}), (r:Region {name: $regionName})
	MERGE (c)-[:BELONGS]->(r)
	`

	_, err := session.ExecuteWrite(ctx, neo4j.ManagedTransactionWork(func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, query, map[string]interface{}{
			"countryCode": strings.ToUpper(strings.TrimSpace(countryCode)),
			"regionName":  strings.TrimSpace(regionName),
		})
		return nil, err
	}))

	return err
}

func (r *Neo4jCountryRepository) AddRegionToCountry(ctx context.Context, countryCode string, region *entities.Region) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close(ctx)

	query := `
	MATCH (c:Country {code: $countryCode})
	CREATE (r:Region {name: $regionName})
	MERGE (c)-[:BELONGS]->(r)
	`

	_, err := session.ExecuteWrite(ctx, neo4j.ManagedTransactionWork(func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, query, map[string]interface{}{
			"countryCode": strings.ToUpper(strings.TrimSpace(countryCode)),
			"regionName":  strings.TrimSpace(region.Name),
		})
		return nil, err
	}))

	return err
}
