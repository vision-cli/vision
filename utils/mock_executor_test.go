package utils_test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vision-cli/vision/utils"
)

func TestCommandExists_ReturnsTrueIfCommandExists(t *testing.T) {
	e := utils.NewMockExecutor()
	e.AddCommand("command")
	assert.True(t, e.CommandExists("command"))
}

func TestCommandExists_ReturnsFalseIfCommandDoesntExist(t *testing.T) {
	e := utils.NewMockExecutor()
	e.AddCommand("command")
	assert.False(t, e.CommandExists("command1"))
}

func TestHistory_ReturnErrorsAndOutput(t *testing.T) {
	e := utils.NewMockExecutor()
	e.Errors(&exec.Cmd{}, "", "error")
	e.Output(&exec.Cmd{}, "", "output")
	assert.Equal(t, []string{"error", "output"}, e.History())
}
