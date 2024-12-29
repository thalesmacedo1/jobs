package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/thalesmacedo1/covid-api/domain/entities"
	"github.com/thalesmacedo1/covid-api/domain/repositories"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jVaccineRepository struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jVaccineRepository(driver neo4j.DriverWithContext) repositories.VaccineRepository {
	return &Neo4jVaccineRepository{
		driver: driver,
	}
}

func (r *Neo4jVaccineRepository) GetVaccinesUsed(ctx context.Context, countryCode string) ([]struct {
	Vaccine           entities.Vaccine
	StartDate         time.Time
	AuthorizationDate time.Time
}, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	query := `
	MATCH (c:Country {code: $countryCode})-[:USES]->(v:Vaccine)
	OPTIONAL MATCH (v)-[:STARTED_ON]->(d:Date)
	OPTIONAL MATCH (v)-[:AUTHORIZATION_ON]->(ad:Date)
	RETURN v.company AS company, v.vaccine AS vaccine, v.product AS product, d.date AS startDate, ad.date AS authorizationDate
	`

	result, err := session.ExecuteRead(ctx, neo4j.ManagedTransactionWork(func(tx neo4j.ManagedTransaction) (interface{}, error) {
		rec, err := tx.Run(ctx, query, map[string]interface{}{
			"countryCode": countryCode,
		})
		if err != nil {
			return nil, err
		}

		var vaccines []struct {
			Vaccine           entities.Vaccine
			StartDate         time.Time
			AuthorizationDate time.Time
		}

		for rec.Next(ctx) {
			record := rec.Record()

			companyInterface, _ := record.Get("company")
			vaccineNameInterface, _ := record.Get("vaccine")
			productInterface, _ := record.Get("product")
			startDateStr, _ := record.Get("startDate")
			authorizationDateStr, _ := record.Get("authorizationDate")

			// Check and convert company
			company, ok := companyInterface.(string)
			if !ok || company == "" {
				company = "Unknown Company"
			}

			// Check and convert vaccine name
			vaccineName, ok := vaccineNameInterface.(string)
			if !ok || vaccineName == "" {
				continue
			}

			// Check and convert product
			product, ok := productInterface.(string)
			if !ok || product == "" {
				product = "Unknown Product"
			}

			// Parse start date
			var startDate time.Time
			if startDateStr != nil {
				if dateStr, ok := startDateStr.(string); ok {
					parsedDate, err := time.Parse("2006-01-02", dateStr)
					if err == nil {
						startDate = parsedDate
					}
				}
			}

			// Parse authorization date
			var authorizationDate time.Time
			if authorizationDateStr != nil {
				if dateStr, ok := authorizationDateStr.(string); ok {
					parsedDate, err := time.Parse("2006-01-02", dateStr)
					if err == nil {
						authorizationDate = parsedDate
					}
				}
			}

			vaccine := entities.NewVaccine(vaccineName, company, product)

			vaccines = append(vaccines, struct {
				Vaccine           entities.Vaccine
				StartDate         time.Time
				AuthorizationDate time.Time
			}{
				Vaccine:           *vaccine,
				StartDate:         startDate,
				AuthorizationDate: authorizationDate,
			})
		}

		if err = rec.Err(); err != nil {
			return nil, err
		}

		return vaccines, nil
	}))

	if err != nil {
		return nil, err
	}

	return result.([]struct {
		Vaccine           entities.Vaccine
		StartDate         time.Time
		AuthorizationDate time.Time
	}), nil
}

func (r *Neo4jVaccineRepository) GetMostUsedVaccine(ctx context.Context, regionName string) (*entities.Vaccine, int, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	query := `
	MATCH (r:Region {name: $regionName})<-[:BELONGS]-(c:Country)-[:USES]->(v:Vaccine)
	RETURN v.name AS vaccineName, v.company AS company, v.product AS product, COUNT(v) AS usageCount
	ORDER BY usageCount DESC
	LIMIT 1
	`

	result, err := session.ExecuteRead(ctx, neo4j.ManagedTransactionWork(func(tx neo4j.ManagedTransaction) (interface{}, error) {
		rec, err := tx.Run(ctx, query, map[string]interface{}{
			"regionName": regionName,
		})
		if err != nil {
			return nil, err
		}

		if rec.Next(ctx) {
			record := rec.Record()
			vaccineNameInterface, _ := record.Get("vaccineName")
			companyInterface, _ := record.Get("company")
			productInterface, _ := record.Get("product")
			usageCountInterface, _ := record.Get("usageCount")

			// Check and convert vaccineName
			vaccineName, ok := vaccineNameInterface.(string)
			if !ok || vaccineName == "" {
				return nil, fmt.Errorf("vaccine name is missing or invalid")
			}

			// Check and convert company
			company, ok := companyInterface.(string)
			if !ok || company == "" {
				company = "Unknown Company"
			}

			// Check and convert product
			product, ok := productInterface.(string)
			if !ok || product == "" {
				product = "Unknown Product"
			}

			// Convert usageCount
			usageCount, ok := usageCountInterface.(int64)
			if !ok {
				return nil, fmt.Errorf("usage count is missing or invalid")
			}

			vaccine := entities.NewVaccine(vaccineName, company, product)

			return struct {
				Vaccine    entities.Vaccine
				UsageCount int
			}{
				Vaccine:    *vaccine,
				UsageCount: int(usageCount),
			}, nil
		}

		if err = rec.Err(); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("no vaccines found for region %s", regionName)
	}))

	if err != nil {
		return nil, 0, err
	}

	res := result.(struct {
		Vaccine    entities.Vaccine
		UsageCount int
	})

	return &res.Vaccine, res.UsageCount, nil
}
