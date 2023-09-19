package transpiler

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/vision-cli/vision/common/transpiler/model"
)

func Load(filename string) model.Project {
	var project model.Project

	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(yamlFile, &project); err != nil {
		panic(err)
	}

	if ok, err := project.IsValid(); !ok {
		panic(err)
	}

	return project
}
