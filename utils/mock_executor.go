package utils

import (
	"os/exec"
)

type MockExecutor struct {
	history []string
	cmds    map[string]string
	output  string
}

func (e *MockExecutor) Errors(cmd *exec.Cmd, targetDir string, action string) error {
	e.history = append(e.history, action)
	return nil
}

func (e *MockExecutor) Output(cmd *exec.Cmd, targetDir string, action string) (string, error) {
	e.history = append(e.history, action)
	return e.output, nil
}

func (e *MockExecutor) CommandExists(cmd string) bool {
	_, exists := e.cmds[cmd]
	return exists
}

func (e *MockExecutor) AddCommand(cmd string) {
	e.cmds[cmd] = cmd
}

func (e *MockExecutor) SetOutput(o string) {
	e.output = o
}

func (e *MockExecutor) History() []string {
	return e.history
}

func NewMockExecutor() *MockExecutor {
	return &MockExecutor{
		history: []string{},
		cmds:    map[string]string{},
	}
}
