package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vision-cli/vision/utils"
)

func TestCheckTools_MustFindGo(t *testing.T) {
	e := utils.NewMockExecutor()
	err := checkTools(e)
	require.Error(t, err)
}

func TestCheckTools_PassIfFindGo(t *testing.T) {
	e := utils.NewMockExecutor()
	e.AddCommand("go")
	err := checkTools(e)
	require.NoError(t, err)
}
