package model

import (
	"errors"
	"fmt"
)

type Module struct {
	ApiVersion string    `yaml:"apiVersion"`
	Name       string    `yaml:"name"`
	Services   []Service `yaml:"services"`
}

func (m Module) isValid() (bool, error) {
	if m.Name == "" {
		return false, errors.New("Module name is required")
	}

	if m.ApiVersion == "" {
		return false, errors.New("Module version is empty")
	}

	if len(m.Services) == 0 {
		return false, fmt.Errorf("Module %s must have at least one service", m.Name)
	}

	for _, service := range m.Services {
		if ok, err := service.isValid(); !ok {
			return false, err
		}
	}

	return true, nil
}
