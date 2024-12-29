package usecases

import (
	"context"
	"time"

	"github.com/thalesmacedo1/covid-api/domain/repositories"
)

type GetCovidTotalsUseCase interface {
	Execute(ctx context.Context, input GetCovidTotalsInput) (*GetCovidTotalsOutput, error)
}

type GetCovidTotalsInput struct {
	CountryCode string
	Date        time.Time
}

type GetCovidTotalsOutput struct {
	CumulativeCases  int
	NewCases         int
	CumulativeDeaths int
	NewDeaths        int
}

type getCovidTotalsUseCase struct {
	covidStatsRepo repositories.CovidStatsRepository
}

func NewGetCovidTotalsUseCase(repo repositories.CovidStatsRepository) GetCovidTotalsUseCase {
	return &getCovidTotalsUseCase{
		covidStatsRepo: repo,
	}
}

func (uc *getCovidTotalsUseCase) Execute(ctx context.Context, input GetCovidTotalsInput) (*GetCovidTotalsOutput, error) {
	stats, err := uc.covidStatsRepo.GetTotalCasesAndDeaths(ctx, input.CountryCode, input.Date)
	if err != nil {
		return nil, err
	}

	return &GetCovidTotalsOutput{
		CumulativeCases:  stats.CumulativeCases,
		NewCases:         stats.NewCases,
		CumulativeDeaths: stats.CumulativeDeaths,
		NewDeaths:        stats.NewDeaths,
	}, nil
}
