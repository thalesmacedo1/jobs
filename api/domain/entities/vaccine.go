package entities

type Vaccine struct {
	Product string
	Company string
	Vaccine string
}

func NewVaccine(
	product string,
	company string,
	vax string,
) *Vaccine {
	vaccine := Vaccine{
		Product: product,
		Company: company,
		Vaccine: vax,
	}
	return &vaccine
}
