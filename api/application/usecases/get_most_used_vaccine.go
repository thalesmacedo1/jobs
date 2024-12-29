package usecases

import (
	"context"

	"github.com/thalesmacedo1/covid-api/domain/entities"
	"github.com/thalesmacedo1/covid-api/domain/repositories"
)

type GetMostUsedVaccineUseCase interface {
	Execute(ctx context.Context, input GetMostUsedVaccineInput) (*GetMostUsedVaccineOutput, error)
}

type GetMostUsedVaccineInput struct {
	RegionName string
}

type GetMostUsedVaccineOutput struct {
	Vaccine entities.Vaccine
	Usage   int
}

type getMostUsedVaccineUseCase struct {
	vaccineRepo repositories.VaccineRepository
}

func NewGetMostUsedVaccineUseCase(repo repositories.VaccineRepository) GetMostUsedVaccineUseCase {
	return &getMostUsedVaccineUseCase{
		vaccineRepo: repo,
	}
}

func (uc *getMostUsedVaccineUseCase) Execute(ctx context.Context, input GetMostUsedVaccineInput) (*GetMostUsedVaccineOutput, error) {
	mostUsedVaccine, usage, err := uc.vaccineRepo.GetMostUsedVaccine(ctx, input.RegionName)
	if err != nil {
		return nil, err
	}

	return &GetMostUsedVaccineOutput{
		Vaccine: *mostUsedVaccine,
		Usage:   usage,
	}, nil
}
