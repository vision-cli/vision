package model

import (
	"errors"
	"fmt"
	"strings"

	"github.com/vision-cli/vision/common/cases"
)

const (
	nameSize = 50
)

type Service struct {
	Name     string   `yaml:"service"`
	Enums    []Enum   `yaml:"enums"`
	Entities []Entity `yaml:"entities"`
}

func (s *Service) HasPersistence() bool {
	for _, e := range s.Entities {
		if e.Persistence == PersistenceDb {
			return true
		}
	}
	return false
}

func (s *Service) HasTimestamp() bool {
	for _, e := range s.Entities {
		for _, f := range e.Fields {
			if f.Type == TypeTimestamp {
				return true
			}
		}
	}
	return false
}

func (s *Service) isValid() (bool, error) {
	if s.Name == "" {
		return false, errors.New("Service name is required")
	}

	if err := NameErr(s.Name); err != nil {
		return false, err
	}

	if len(s.Entities) == 0 {
		return false, errors.New("Service must have at least one entity")
	}

	for _, entity := range s.Entities {
		if ok, err := entity.isValid(); !ok {
			return false, err
		}
	}

	return true, nil
}

// NameErr reports errors on names conflicting with k8s and vision naming conventions.
func NameErr(name string) error {
	if kebab := cases.Kebab(name); kebab != name {
		return fmt.Errorf("must be kebab-case (suggestion: %q)", kebab)
	}

	banned := map[string]struct{}{
		"all":       {},
		"create":    {},
		"get":       {},
		"namespace": {},
		"service":   {},
		"srv":       {},
		"svc":       {},
	}

	for _, word := range strings.Split(name, "-") {
		if _, match := banned[word]; match {
			return fmt.Errorf("must not contain the word %s", word)
		}
	}

	if len(name) > nameSize {
		return fmt.Errorf("must not be longer than %d characters", nameSize)
	}

	return nil
}
