package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vision-cli/vision/execute"
)

func TestCheckTools_MustFindGo(t *testing.T) {
	e := execute.NewMockExecutor()
	err := checkTools(e)
	require.Error(t, err)
}

func TestCheckTools_PassIfFindGo(t *testing.T) {
	e := execute.NewMockExecutor()
	e.AddCommand("go")
	err := checkTools(e)
	require.NoError(t, err)
}
