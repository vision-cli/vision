package model

type Enum struct {
	Name   string   `yaml:"enum"`
	Values []string `yaml:"values"`
}
