package entities

import (
	"errors"
	"strings"
)

type Region struct {
	Name string
}

func NewRegion(name string) (*Region, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return nil, errors.New("region name cannot be empty")
	}

	return &Region{
		Name: name,
	}, nil
}
