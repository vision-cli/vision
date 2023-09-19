package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vision-cli/vision/common/mocks"
)

func TestCheckTools_MustFindGo(t *testing.T) {
	e := mocks.NewMockExecutor()
	err := checkTools(&e)
	require.Error(t, err)
}

func TestCheckTools_PassIfFindGo(t *testing.T) {
	e := mocks.NewMockExecutor()
	e.AddCommand("go")
	err := checkTools(&e)
	require.NoError(t, err)
}
