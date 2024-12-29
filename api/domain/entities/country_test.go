package entities

import (
	"testing"
)

func TestNewCountry(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		country  string
		wantErr  bool
		errMsg   string
		wantCode string
		wantName string
	}{
		{
			name:     "Valid 2-letter code",
			code:     "us",
			country:  "United States",
			wantErr:  false,
			wantCode: "US",
			wantName: "United States",
		},
		{
			name:     "Valid 3-letter code with spaces",
			code:     "  bra ",
			country:  "Brazil ",
			wantErr:  false,
			wantCode: "BRA",
			wantName: "Brazil",
		},
		{
			name:    "Empty country code",
			code:    " ",
			country: "Canada",
			wantErr: true,
			errMsg:  "country code cannot be empty",
		},
		{
			name:    "Invalid country code length",
			code:    "usa1",
			country: "USA",
			wantErr: true,
			errMsg:  "country code must be 2 or 3 characters long",
		},
		{
			name:    "Empty country name",
			code:    "CA",
			country: " ",
			wantErr: true,
			errMsg:  "country name cannot be empty",
		},
		{
			name:     "Code with mixed case and spaces",
			code:     "  cA ",
			country:  " Canada ",
			wantErr:  false,
			wantCode: "CA",
			wantName: "Canada",
		},
		{
			name:     "3-letter code with valid name",
			code:     "ind",
			country:  "India",
			wantErr:  false,
			wantCode: "IND",
			wantName: "India",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			country, err := NewCountry(tt.code, tt.country)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCountry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err.Error() != tt.errMsg {
					t.Errorf("NewCountry() error message = '%v', want '%v'", err.Error(), tt.errMsg)
				}
				return
			}
			if country.Code != tt.wantCode {
				t.Errorf("NewCountry() Code = %v, want %v", country.Code, tt.wantCode)
			}
			if country.Name != tt.wantName {
				t.Errorf("NewCountry() Name = %v, want %v", country.Name, tt.wantName)
			}
		})
	}
}
