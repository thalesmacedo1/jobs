package usecases

import (
	"context"

	"github.com/thalesmacedo1/covid-api/domain/entities"
	"github.com/thalesmacedo1/covid-api/domain/repositories"
)

type GetVaccinesUsedUseCase interface {
	Execute(ctx context.Context, input GetVaccinesUsedInput) ([]GetVaccinesUsedOutput, error)
}

type GetVaccinesUsedInput struct {
	CountryCode string
}

type GetVaccinesUsedOutput struct {
	Vaccine           entities.Vaccine `json:"vaccine"`
	StartDate         string           `json:"start_date"`
	AuthorizationDate string           `json:"authorization_date,omitempty"`
}

type getVaccinesUsedUseCase struct {
	vaccineRepo repositories.VaccineRepository
}

func NewGetVaccinesUsedUseCase(repo repositories.VaccineRepository) GetVaccinesUsedUseCase {
	return &getVaccinesUsedUseCase{
		vaccineRepo: repo,
	}
}

func (uc *getVaccinesUsedUseCase) Execute(ctx context.Context, input GetVaccinesUsedInput) ([]GetVaccinesUsedOutput, error) {
	vaccinesUsed, err := uc.vaccineRepo.GetVaccinesUsed(ctx, input.CountryCode)
	if err != nil {
		return nil, err
	}

	uniqueVaccines := make(map[string]bool)
	var output []GetVaccinesUsedOutput

	for _, vu := range vaccinesUsed {

		vaccineName := vu.Vaccine.Vaccine

		if _, exists := uniqueVaccines[vaccineName]; exists {
			continue
		}

		uniqueVaccines[vaccineName] = true

		vaccineOutput := GetVaccinesUsedOutput{
			Vaccine:   vu.Vaccine,
			StartDate: vu.StartDate.Format("2006-01-02"),
		}

		if !vu.AuthorizationDate.IsZero() {
			vaccineOutput.AuthorizationDate = vu.AuthorizationDate.Format("2006-01-02")
		}

		output = append(output, vaccineOutput)
	}

	return output, nil
}
