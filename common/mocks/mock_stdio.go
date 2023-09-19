package mocks

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func WithMockStdio(t *testing.T, input string, test func()) string {
	t.Helper()
	return WithMockStdout(t, func() {
		WithMockStdin(t, input, test)
	})
}

func WithMockStdout(t *testing.T, test func()) string {
	t.Helper()
	stdoutMock, err := os.CreateTemp("", "stdoutMock-*")
	require.NoError(t, err)
	defer os.Remove(stdoutMock.Name())

	oldStdout := os.Stdout
	os.Stdout = stdoutMock
	defer func() {
		os.Stdout = oldStdout
	}()

	test()

	_, err = stdoutMock.Seek(0, 0)
	require.NoError(t, err)

	output, err := io.ReadAll(stdoutMock)
	require.NoError(t, err)
	return string(output)
}

func WithMockStdin(t *testing.T, input string, test func()) {
	t.Helper()
	if input[len(input)-1] != '\n' {
		input += "\n"
	}
	stdinMock, err := os.CreateTemp("", "stdinMock-*")
	require.NoError(t, err)
	defer os.Remove(stdinMock.Name())

	_, err = stdinMock.Write([]byte(input))
	require.NoError(t, err)

	_, err = stdinMock.Seek(0, 0)
	require.NoError(t, err)

	oldStdin := os.Stdin
	os.Stdin = stdinMock
	defer func() {
		os.Stdin = oldStdin
	}()

	test()
}
