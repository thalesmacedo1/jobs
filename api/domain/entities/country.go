package entities

import (
	"errors"
	"strings"
)

type Country struct {
	Code string
	Name string
}

func NewCountry(code, name string) (*Country, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	name = strings.TrimSpace(name)

	if code == "" {
		return nil, errors.New("country code cannot be empty")
	}

	if len(code) != 2 && len(code) != 3 {
		return nil, errors.New("country code must be 2 or 3 characters long")
	}

	if name == "" {
		return nil, errors.New("country name cannot be empty")
	}

	return &Country{
		Code: code,
		Name: name,
	}, nil
}
