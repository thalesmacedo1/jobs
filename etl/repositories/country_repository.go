package repositories

import (
	"context"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type CountryRepository struct {
	driver neo4j.DriverWithContext
}

func NewCountryRepository(driver neo4j.DriverWithContext) *CountryRepository {
	return &CountryRepository{driver: driver}
}

func (r *CountryRepository) CreateCountry(ctx context.Context, country, code, region string) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.Run(ctx,
		`MERGE (c:Country {code: $code})
         SET c.name = $name
         MERGE (r:Region {name: $region})
         MERGE (c)-[:BELONGS]->(r)`,
		map[string]interface{}{
			"name":   country,
			"code":   code,
			"region": region,
		},
	)
	if err != nil {
		log.Printf("Error creating country %s: %v", country, err)
	}
	return err
}
