package model

import (
	"errors"
)

type Project struct {
	Name                    string   `yaml:"project"`
	Deployment              string   `yaml:"deployment"`
	Remote                  string   `yaml:"remote"`
	DefaultServiceNamespace string   `yaml:"default-namespace"`
	Registry                string   `yaml:"registry"`
	Branch                  string   `yaml:"branch"`
	Modules                 []Module `yaml:"modules"`
}

func (p Project) IsValid() (bool, error) {
	// check name and modules only
	// the remaining project fields, although required, can be provided through flags
	if p.Name == "" {
		return false, errors.New("Project name is required")
	}

	if len(p.Modules) == 0 {
		return false, errors.New("Project must have at least one module")
	}

	for _, module := range p.Modules {
		if ok, err := module.isValid(); !ok {
			return false, err
		}
	}

	return true, nil
}
