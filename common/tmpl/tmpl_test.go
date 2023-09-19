package tmpl_test

import (
	"embed"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vision-cli/vision/common/file"
	"github.com/vision-cli/vision/common/mocks"
	"github.com/vision-cli/vision/common/tmpl"
)

//go:embed all:_templates
var templateFiles embed.FS

func TestIsTemplate_ForTemplates_ReturnsTrue(t *testing.T) {
	assert.True(t, tmpl.IsTemplate("file.tmpl"))
}

func TestIsTemplate_ForNonTemplates_ReturnsFalse(t *testing.T) {
	assert.False(t, tmpl.IsTemplate("file.go"))
}

func TestGenerateFS_NoSkipExisting_GeneratesCorrectFS(t *testing.T) {
	tw := mocks.NewMockTmplWriter()
	err := tmpl.GenerateFS(templateFiles, "_templates", "out", nil, false, &tw)
	require.NoError(t, err)
	assert.Equal(t, []string{"CreateDir: out", "WriteTemplatedFS: out/file", "WriteExactFS: out/random.file"}, tw.History())
}

func TestGenerateFS_SkipExisting_GeneratesCorrectFS(t *testing.T) {
	old := file.Osstat
	defer func() { file.Osstat = old }()
	file.Osstat = func(name string) (fs.FileInfo, error) {
		return nil, nil
	}

	tw := mocks.NewMockTmplWriter()
	err := tmpl.GenerateFS(templateFiles, "_templates", "out", nil, true, &tw)
	require.NoError(t, err)
	assert.Equal(t, []string{"CreateDir: out"}, tw.History())
}
