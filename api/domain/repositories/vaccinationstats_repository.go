package repositories

import (
	"context"
	"time"

	"github.com/thalesmacedo1/covid-api/domain/valueobjects"
)

type VaccinationStatsRepository interface {
	GetVaccinatedPeople(ctx context.Context, countryCode string, date time.Time) (*valueobjects.VaccinationStats, error)
}
