package plugin

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type Executor struct {
	FullPath string
	Args     []string
}

func NewExecutor(path string, args []string) Executor {
	return Executor{
		FullPath: path,
		Args:     args,
	}
}

type Info struct {
	ShortDescription string `json:"short_description"`
	LongDescription  string `json:"long_description"`
	Flags            []*flag
}

type flag struct {
	Name      string
	Shorthand string
	Usage     string
	Type      string
}

// info returns usage and descriptions of the plugin
// TODO(steve): make info resp part of the plugin API
func (e Executor) Info() (*Info, error) {
	cmd := exec.Command(e.FullPath, "info")
	b, err := cmd.Output()
	if err != nil {
		return nil, err
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
	// hash??
	// git sha ??
}

// TODO(steve): make version resp part of the plugin API
func (e Executor) Version() (*Version, error) {
	cmd := exec.Command(e.FullPath, "version")
	b, err := cmd.Output()
	if err != nil {
		return nil, err
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
	// if e.Args[0] != "init" {
	// 	return nil, fmt.Errorf("Expected init to be called")
	// }
	cmd := exec.Command(e.FullPath, e.Args...)
	b, err := cmd.Output()
	if err != nil {
		return nil, err
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
		return nil, err
	}
	var g Generate
	err = json.Unmarshal(b, &g)
	if err != nil {
		return nil, err
	}
	return &g, nil
}
