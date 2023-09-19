package mocks_test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vision-cli/vision/common/mocks"
)

func TestCommandExists_ReturnsTrueIfCommandExists(t *testing.T) {
	e := mocks.NewMockExecutor()
	e.AddCommand("command")
	assert.True(t, e.CommandExists("command"))
}

func TestCommandExists_ReturnsFalseIfCommandDoesntExist(t *testing.T) {
	e := mocks.NewMockExecutor()
	e.AddCommand("command")
	assert.False(t, e.CommandExists("command1"))
}

func TestHistory_ReturnErrorsAndOutput(t *testing.T) {
	e := mocks.NewMockExecutor()
	err := e.Errors(&exec.Cmd{}, "", "error")
	require.NoError(t, err)
	_, err = e.Output(&exec.Cmd{}, "", "output")
	require.NoError(t, err)
	assert.Equal(t, []string{"error", "output"}, e.History())
}

func TestOutput_ReturnSetOutput(t *testing.T) {
	e := mocks.NewMockExecutor()
	e.SetOutput("output")
	r, err := e.Output(&exec.Cmd{}, "", "output")
	assert.NoError(t, err)
	assert.Equal(t, "output", r)
}
