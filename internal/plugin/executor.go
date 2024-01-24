package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/log"
	api "github.com/vision-cli/api/v1"
)

type Executor struct {
	FullPath string
}

func NewExecutor(path string) Executor {
	return Executor{
		FullPath: path,
	}
}

// info returns usage and descriptions of the plugin
// TODO(steve): make info resp part of the plugin API
func (e Executor) Info() (*api.Info, error) {
	cmd := exec.Command(e.FullPath, "info")
	b, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("info command: %w", err)
	}
	var i api.Info
	err = json.Unmarshal(b, &i)
	if err != nil {
		return nil, fmt.Errorf("info: invalid json resp from plugin: %w", err)
	}
	return &i, nil
}

func (e Executor) Version() (*api.Version, error) {
	cmd := exec.Command(e.FullPath, "version")
	b, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("version command: %w", err)
	}
	var v api.Version
	err = json.Unmarshal(b, &v)
	if err != nil {
		return nil, fmt.Errorf("version: invalid json resp from plugin: %w", err)
	}
	return &v, nil
}

func (e Executor) Init() (*api.Init, error) {
	var stderrBuf strings.Builder
	cmd := exec.Command(e.FullPath, os.Args[2:]...)
	cmd.Stderr = &stderrBuf

	b, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("init command: %w", err)
	}

	stderrStr := stderrBuf.String()
	if stderrStr != "" {
		fmt.Println(stderrStr)
		return nil, fmt.Errorf("init command failed")
	}

	var i api.Init

	err = json.Unmarshal(b, &i)
	if err != nil {
		return nil, fmt.Errorf("init: invalid json resp from plugin: %w", err)
	}

	if !i.Success {
		return nil, fmt.Errorf("init command failed")
	}

	log.Info("Plugin init successful")
	return &i, nil
}

func (e Executor) Generate() (*api.Generate, error) {
	var stderrBuf strings.Builder
	cmd := exec.Command(e.FullPath, "generate")
	cmd.Stderr = &stderrBuf

	b, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("generate command: %w", err)
	}

	stderrStr := stderrBuf.String()
	if stderrStr != "" {
		fmt.Println(stderrStr)
		return nil, fmt.Errorf("init command failed")
	}

	var g api.Generate
	err = json.Unmarshal(b, &g)
	if err != nil {
		return nil, fmt.Errorf("generate: invalid json resp from plugin: %w", err)
	}

	if !g.Success {
		return nil, fmt.Errorf("init command failed")
	}

	return &g, nil
}
