package workspace_test

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vision-cli/vision/common/file"
	"github.com/vision-cli/vision/common/mocks"
	"github.com/vision-cli/vision/common/workspace"
)

func TestInit_RunsGoWorkInit(t *testing.T) {
	old := file.Osremoveall
	defer func() { file.Osremoveall = old }()
	removed := []string{}
	file.Osremoveall = func(path string) error {
		removed = append(removed, path)
		return nil
	}
	e := mocks.NewMockExecutor()
	result := workspace.Init("targetdir", &e)
	require.NoError(t, result)
	assert.Equal(t, 2, len(removed))
	assert.Equal(t, "targetdir/go.work", removed[0])
	assert.Equal(t, "targetdir/go.work.sum", removed[1])
	assert.Equal(t, 1, len(e.History()))
	assert.Equal(t, "initialising workspace", e.History()[0])
}

func TestInit_FileCannotBeRemoved_ReturnsError(t *testing.T) {
	old := file.Osremoveall
	defer func() { file.Osremoveall = old }()
	file.Osremoveall = func(path string) error {
		return errors.New("cannot remove file")
	}
	e := mocks.NewMockExecutor()
	result := workspace.Init("targetdir", &e)
	require.Error(t, result)
}

func TestUse_WorkspaceExists_RunsGoWorkUsePath(t *testing.T) {
	old := file.Osstat
	defer func() { file.Osstat = old }()
	file.Osstat = func(name string) (os.FileInfo, error) {
		return nil, nil
	}
	e := mocks.NewMockExecutor()
	result := workspace.Use("targetdir", "targetPath", &e)
	require.NoError(t, result)
	assert.Equal(t, 1, len(e.History()))
	assert.Equal(t, "updating workspace modules", e.History()[0])
}

func TestUse_WorkspaceDoesNotExist_RunsInitAndGoWorkUsePath(t *testing.T) {
	oldosstat := file.Osstat
	defer func() { file.Osstat = oldosstat }()
	file.Osstat = func(name string) (os.FileInfo, error) {
		return nil, os.ErrNotExist
	}
	oldremoveall := file.Osremoveall
	defer func() { file.Osremoveall = oldremoveall }()
	file.Osremoveall = func(path string) error {
		return nil
	}
	e := mocks.NewMockExecutor()
	result := workspace.Use("targetdir", "targetPath", &e)
	require.NoError(t, result)
	assert.Equal(t, 2, len(e.History()))
	assert.Equal(t, "initialising workspace", e.History()[0])
	assert.Equal(t, "updating workspace modules", e.History()[1])
}
