package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type Executor struct {
	FullPath string
}

func NewExecutor(path string) Executor {
	return Executor{
		FullPath: path,
	}
}

type Info struct {
	ShortDescription string `json:"short_description"`
	LongDescription  string `json:"long_description"`
}

// info returns usage and descriptions of the plugin
// TODO(steve): make info resp part of the plugin API
func (e Executor) Info() (*Info, error) {
	cmd := exec.Command(e.FullPath, "info")
	b, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("info command: %w", err)
	}
	var i Info
	err = json.Unmarshal(b, &i)
	if err != nil {
		return nil, fmt.Errorf("info: invalid json resp from plugin: %w", err)
	}
	return &i, nil
}

type Version struct {
	SemVer string `json:"sem_ver"`
}

func (e Executor) Version() (*Version, error) {
	cmd := exec.Command(e.FullPath, "version")
	b, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("version command: %w", err)
	}
	var v Version
	err = json.Unmarshal(b, &v)
	if err != nil {
		return nil, fmt.Errorf("version: invalid json resp from plugin: %w", err)
	}
	return &v, nil
}

type Init struct {
	Config any `json:"config"`
}

func (e Executor) Init() (*Init, error) {
	var loc string
	if os.Args[1] == "" {
		loc = ""
	} else {
		loc = os.Args[1]
	}
	cmd := exec.Command(e.FullPath, "init", loc)
	b, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("init command: %w", err)
	}
	var i Init
	err = json.Unmarshal(b, &i)
	if err != nil {
		return nil, fmt.Errorf("init: invalid json resp from plugin: %w", err)
	}
	return &i, nil
}

type Generate struct {
	Success bool `json:"success"`
}

func (e Executor) Generate() (*Generate, error) {
	cmd := exec.Command(e.FullPath, "generate")
	b, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("generate command: %w", err)
	}
	var g Generate
	err = json.Unmarshal(b, &g)
	if err != nil {
		return nil, fmt.Errorf("generate: invalid json resp from plugin: %w", err)
	}
	return &g, nil
}
