package entities

import (
	"testing"
)

func TestNewRegion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		errMsg   string
		wantName string
	}{
		{
			name:     "Valid region name",
			input:    "Europe",
			wantErr:  false,
			wantName: "Europe",
		},
		{
			name:     "Region name with leading and trailing spaces",
			input:    "  Asia ",
			wantErr:  false,
			wantName: "Asia",
		},
		{
			name:    "Empty region name",
			input:   "   ",
			wantErr: true,
			errMsg:  "region name cannot be empty",
		},
		{
			name:     "Region name with mixed case and spaces",
			input:    " North America ",
			wantErr:  false,
			wantName: "North America",
		},
		{
			name:     "Single character region name",
			input:    "A",
			wantErr:  false,
			wantName: "A",
		},
		{
			name:     "Region name with special characters",
			input:    "Oceania!",
			wantErr:  false,
			wantName: "Oceania!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			region, err := NewRegion(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRegion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err.Error() != tt.errMsg {
					t.Errorf("NewRegion() error message = '%v', want '%v'", err.Error(), tt.errMsg)
				}
				return
			}
			if region.Name != tt.wantName {
				t.Errorf("NewRegion() Name = %v, want %v", region.Name, tt.wantName)
			}
		})
	}
}
