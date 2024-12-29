package repositories

import (
	"context"
	"time"

	"github.com/thalesmacedo1/covid-api/domain/entities"
)

type VaccineRepository interface {
	GetVaccinesUsed(ctx context.Context, countryCode string) ([]struct {
		Vaccine           entities.Vaccine
		StartDate         time.Time
		AuthorizationDate time.Time
	}, error)

	GetMostUsedVaccine(ctx context.Context, regionName string) (*entities.Vaccine, int, error)
}
