package plugin

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/charmbracelet/log"
)

type Executor struct {
	FullPath string
}

func NewExecutor(path string) Executor {
	return Executor{FullPath: path}
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
	cmd := exec.Command(e.FullPath, "init")
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

// Executes a command and an argument to a Vision plugin binary then processes any output.
func (e Executor) RunCommand(pluginName string, arg string) error {
	root := fmt.Sprintf("vision-plugin-%v-v1", pluginName)
	cmd := exec.Command(root, arg)
	bts, err := cmd.Output()
	if err != nil {
		return err
	}
	log.Info(string(bts))
	return nil
}