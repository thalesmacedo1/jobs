package usecases

import (
	"context"
	"time"

	"github.com/thalesmacedo1/covid-api/domain/repositories"
)

type GetVaccinatedPeopleUseCase interface {
	Execute(ctx context.Context, input GetVaccinatedPeopleInput) (*GetVaccinatedPeopleOutput, error)
}

type GetVaccinatedPeopleInput struct {
	CountryCode string
	Date        time.Time
}

type GetVaccinatedPeopleOutput struct {
	PersonsVaccinated1PlusDose       int
	PersonsVaccinated1PlusDosePer100 int
	TotalVaccinations                int
	TotalVaccinationsPer100          int
}

type getVaccinatedPeopleUseCase struct {
	vaccinationStatsRepo repositories.VaccinationStatsRepository
}

func NewGetVaccinatedPeopleUseCase(repo repositories.VaccinationStatsRepository) GetVaccinatedPeopleUseCase {
	return &getVaccinatedPeopleUseCase{
		vaccinationStatsRepo: repo,
	}
}

func (uc *getVaccinatedPeopleUseCase) Execute(ctx context.Context, input GetVaccinatedPeopleInput) (*GetVaccinatedPeopleOutput, error) {
	stats, err := uc.vaccinationStatsRepo.GetVaccinatedPeople(ctx, input.CountryCode, input.Date)
	if err != nil {
		return nil, err
	}

	return &GetVaccinatedPeopleOutput{
		PersonsVaccinated1PlusDose:       stats.PersonsVaccinated1PlusDose,
		PersonsVaccinated1PlusDosePer100: stats.PersonsVaccinated1PlusDosePer100,
		TotalVaccinations:                stats.TotalVaccinations,
		TotalVaccinationsPer100:          stats.TotalVaccinationsPer100,
	}, nil
}
