package repositories

import (
	"context"

	"github.com/thalesmacedo1/covid-api/domain/entities"
)

type CountryRepository interface {
	GetCountryByCode(ctx context.Context, code string) (*entities.Country, error)
}
