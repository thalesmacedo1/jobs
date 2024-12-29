package repositories

import (
	"context"
	"log"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type VaccineRepository struct {
	driver neo4j.DriverWithContext
}

func NewVaccineRepository(driver neo4j.DriverWithContext) *VaccineRepository {
	return &VaccineRepository{driver: driver}
}

func (r *VaccineRepository) CreateVaccine(ctx context.Context, product string, vaccineName string, company string, authorizationDate time.Time, startDate time.Time) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.Run(ctx,
		`MERGE (v:Vaccine {product: $product})
         SET v.vaccine = $vaccineName, v.company = $company
         MERGE (authDate:Date {date: date($authorizationDate)})
         MERGE (startDate:Date {date: date($startDate)})
         MERGE (v)-[:AUTHORIZATION_ON]->(authDate)
         MERGE (v)-[:STARTED_ON]->(startDate)`,
		map[string]interface{}{
			"product":           product,
			"vaccineName":       vaccineName,
			"company":           company,
			"authorizationDate": authorizationDate,
			"startDate":         startDate,
		})
	if err != nil {
		log.Printf("Error creating vaccine %s: %v", product, err)
	}
	return err
}

func (r *VaccineRepository) CreateCountryUsesVaccine(ctx context.Context, countryCode string, product string) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.Run(ctx,
		`MATCH (c:Country {code: $countryCode}), (v:Vaccine {product: $product})
         MERGE (c)-[:USES]->(v)`,
		map[string]interface{}{
			"countryCode": countryCode,
			"product":     product,
		})
	if err != nil {
		log.Printf("Error creating USES relationship between country %s and vaccine %s: %v", countryCode, product, err)
	}
	return err
}
