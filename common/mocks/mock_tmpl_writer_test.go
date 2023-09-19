package mocks_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vision-cli/vision/common/mocks"
)

func TestWriteTemplatedFS_ReturnNilAndAddsToHistory(t *testing.T) {
	tw := mocks.NewMockTmplWriter()
	err := tw.WriteTemplatedFS("templatePath", "targetPath", nil, nil)
	require.NoError(t, err)
	assert.Equal(t, 1, len(tw.History()))
	assert.Equal(t, "WriteTemplatedFS: targetPath", tw.History()[0])
}

func TestWriteExactFS_ReturnNilAndAddsToHistory(t *testing.T) {
	tw := mocks.NewMockTmplWriter()
	err := tw.WriteExactFS("templatePath", "targetPath", nil)
	require.NoError(t, err)
	assert.Equal(t, 1, len(tw.History()))
	assert.Equal(t, "WriteExactFS: targetPath", tw.History()[0])
}

func TestCreateDir_ReturnNilAndAddsToHistory(t *testing.T) {
	tw := mocks.NewMockTmplWriter()
	err := tw.CreateDir("dirname")
	require.NoError(t, err)
	assert.Equal(t, 1, len(tw.History()))
	assert.Equal(t, "CreateDir: dirname", tw.History()[0])
}

func TestExists_ReturnTrueIfExists(t *testing.T) {
	tw := mocks.NewMockTmplWriter()
	tw.AddExists("filename")
	assert.True(t, tw.Exists("filename"))
}

func TestExists_ReturnFalseIfDoesntExists(t *testing.T) {
	tw := mocks.NewMockTmplWriter()
	assert.False(t, tw.Exists("filename"))
}
