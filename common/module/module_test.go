package module_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vision-cli/vision/common/file"
	"github.com/vision-cli/vision/common/mocks"
	"github.com/vision-cli/vision/common/module"
)

func TestInit_RunsGoModInit(t *testing.T) {
	old := file.Osremoveall
	defer func() { file.Osremoveall = old }()
	removed := []string{}
	file.Osremoveall = func(path string) error {
		removed = append(removed, path)
		return nil
	}
	e := mocks.NewMockExecutor()
	result := module.Init("targetdir", "modulename", &e)
	require.NoError(t, result)
	assert.Equal(t, 2, len(removed))
	assert.Equal(t, "targetdir/go.mod", removed[0])
	assert.Equal(t, "targetdir/go.sum", removed[1])
	assert.Equal(t, 1, len(e.History()))
	assert.Equal(t, "initialising module", e.History()[0])
}

func TestTidy_RunsGoModTidy(t *testing.T) {
	e := mocks.NewMockExecutor()
	result := module.Tidy("targetdir", &e)
	require.NoError(t, result)
	assert.Equal(t, "finding required module dependencies", e.History()[0])
}

func TestName_ReturnsModName(t *testing.T) {
	old := file.ToLines
	defer func() { file.ToLines = old }()
	file.ToLines = func(path string) ([]string, error) {
		return []string{"module modulename", "println()", "another line"}, nil
	}
	name, err := module.Name("targetdir")
	require.NoError(t, err)
	assert.Equal(t, "modulename", name)
}

func TestRename_RenamesMod(t *testing.T) {
	oldtolines := file.ToLines
	defer func() { file.ToLines = oldtolines }()
	file.ToLines = func(path string) ([]string, error) {
		return []string{"module modulename", "println()", "another line"}, nil
	}
	oldfromlines := file.FromLines
	defer func() { file.FromLines = oldfromlines }()
	var result []string
	file.FromLines = func(path string, lines []string) error {
		result = lines
		return nil
	}

	err := module.Rename("targetdir", "newmodulename")
	require.NoError(t, err)
	assert.Equal(t, []string{"module newmodulename", "println()", "another line"}, result)
}

func TestReplace_RunsGoModEditReplace(t *testing.T) {
	e := mocks.NewMockExecutor()
	result := module.Replace("targetdir", "servicemod", "replacement", &e)
	require.NoError(t, result)
	assert.Equal(t, "replace for servicemod", e.History()[0])
}
