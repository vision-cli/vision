package tmplutils_test

import (
	"embed"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/atos-digital/NHSS-scigateway/internal/tmplutils"
)

func TestWalkDir(t *testing.T) {
	expected := []string{
		"testdata/path/a.gohtml",
		"testdata/path/b/ab.gohtml",
		"testdata/path/b/c/abc.gohtml",
		"testdata/path/b/c/d/abcd.gohtml",
	}
	paths, err := tmplutils.WalkDir("testdata/path")
	require.Nil(t, err)
	assert.Equal(t, expected, paths)
}

func TestParseDir(t *testing.T) {
	tmpl, err := tmplutils.ParseDir("test", "testdata/path")
	require.Nil(t, err)
	require.NotNil(t, tmpl)
	s := strings.TrimPrefix(tmpl.DefinedTemplates(), "; defined templates are: ")
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\"", "")
	ss := strings.Split(s, ",")
	assert.ElementsMatch(t, []string{"a.gohtml", "ab.gohtml", "abc.gohtml", "abcd.gohtml"}, ss)
}

func TestWalkFS(t *testing.T) {
	expected := []string{
		"a.gohtml",
		"b/ab.gohtml",
		"b/c/abc.gohtml",
		"b/c/d/abcd.gohtml",
	}
	paths, err := tmplutils.WalkFS(os.DirFS("testdata/path"))
	require.Nil(t, err)
	assert.Equal(t, expected, paths)
}

func TestParseFS(t *testing.T) {
	tmpl, err := tmplutils.ParseFS("test", os.DirFS("testdata/path"))
	require.Nil(t, err)
	require.NotNil(t, tmpl)

	s := strings.TrimPrefix(tmpl.DefinedTemplates(), "; defined templates are: ")
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\"", "")
	ss := strings.Split(s, ",")
	assert.ElementsMatch(t, []string{"a.gohtml", "ab.gohtml", "abc.gohtml", "abcd.gohtml"}, ss)
}

//go:embed testdata/path
var testdata embed.FS

func TestParseEmbedFS(t *testing.T) {
	tmpl, err := tmplutils.ParseFS("test", testdata)
	require.Nil(t, err)
	require.NotNil(t, tmpl)

	s := strings.TrimPrefix(tmpl.DefinedTemplates(), "; defined templates are: ")
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\"", "")
	ss := strings.Split(s, ",")
	assert.ElementsMatch(t, []string{"a.gohtml", "ab.gohtml", "abc.gohtml", "abcd.gohtml"}, ss)
}
