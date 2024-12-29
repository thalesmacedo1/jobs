package entities

import (
	"testing"
)

func TestNewVaccine(t *testing.T) {
	tests := []struct {
		name     string
		product  string
		company  string
		vax      string
		wantProd string
		wantCmp  string
		wantVax  string
	}{
		{
			name:     "All fields provided",
			product:  "ProductA",
			company:  "CompanyA",
			vax:      "VaccineA",
			wantProd: "ProductA",
			wantCmp:  "CompanyA",
			wantVax:  "VaccineA",
		},
		{
			name:     "Empty product",
			product:  "",
			company:  "CompanyB",
			vax:      "VaccineB",
			wantProd: "",
			wantCmp:  "CompanyB",
			wantVax:  "VaccineB",
		},
		{
			name:     "Empty company",
			product:  "ProductC",
			company:  "",
			vax:      "VaccineC",
			wantProd: "ProductC",
			wantCmp:  "",
			wantVax:  "VaccineC",
		},
		{
			name:     "Empty vaccine name",
			product:  "ProductD",
			company:  "CompanyD",
			vax:      "",
			wantProd: "ProductD",
			wantCmp:  "CompanyD",
			wantVax:  "",
		},
		{
			name:     "All fields empty",
			product:  "",
			company:  "",
			vax:      "",
			wantProd: "",
			wantCmp:  "",
			wantVax:  "",
		},
		{
			name:     "Fields with spaces",
			product:  "  ProductE  ",
			company:  "  CompanyE  ",
			vax:      "  VaccineE  ",
			wantProd: "  ProductE  ",
			wantCmp:  "  CompanyE  ",
			wantVax:  "  VaccineE  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vaccine := NewVaccine(tt.product, tt.company, tt.vax)

			if vaccine.Product != tt.wantProd {
				t.Errorf("NewVaccine() Product = '%v', want '%v'", vaccine.Product, tt.wantProd)
			}
			if vaccine.Company != tt.wantCmp {
				t.Errorf("NewVaccine() Company = '%v', want '%v'", vaccine.Company, tt.wantCmp)
			}
			if vaccine.Vaccine != tt.wantVax {
				t.Errorf("NewVaccine() Vaccine = '%v', want '%v'", vaccine.Vaccine, tt.wantVax)
			}
		})
	}
}
