package valueobjects

type VaccinationStats struct {
	PersonsBoosterAddDosePer100      int
	PersonsBoosterAddDose            int
	PersonsLastDosePer100            int
	PersonsLastDose                  int
	PersonsVaccinated1PlusDosePer100 int
	TotalVaccinationsPer100          int
	PersonsVaccinated1PlusDose       int
	TotalVaccinations                int
}

func NewVaccinationStats(
	personsBoosterAddDosePer100, personsBoosterAddDose,
	personsLastDosePer100, personsLastDose,
	personsVaccinated1PlusDosePer100, totalVaccinationsPer100,
	personsVaccinated1PlusDose, totalVaccinations int,
) VaccinationStats {
	return VaccinationStats{
		PersonsBoosterAddDosePer100:      personsBoosterAddDosePer100,
		PersonsBoosterAddDose:            personsBoosterAddDose,
		PersonsLastDosePer100:            personsLastDosePer100,
		PersonsLastDose:                  personsLastDose,
		PersonsVaccinated1PlusDosePer100: personsVaccinated1PlusDosePer100,
		TotalVaccinationsPer100:          totalVaccinationsPer100,
		PersonsVaccinated1PlusDose:       personsVaccinated1PlusDose,
		TotalVaccinations:                totalVaccinations,
	}
}
