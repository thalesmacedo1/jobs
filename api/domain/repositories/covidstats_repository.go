package repositories

import (
	"context"
	"time"

	"github.com/thalesmacedo1/covid-api/domain/valueobjects"
)

type CovidStatsRepository interface {
	GetTotalCasesAndDeaths(ctx context.Context, countryCode string, date time.Time) (*valueobjects.CovidStats, error)
	GetCountryWithMostCases(ctx context.Context, date time.Time) (string, int, error)
}
