package utils_test

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vision-cli/vision/utils"
)

func TestWithMockStdIn_ReturnsSingleLine(t *testing.T) {
	utils.WithMockStdin(t, "answer\n", func() {
		reader := bufio.NewReader(os.Stdin)
		answer, err := reader.ReadString('\n')
		require.NoError(t, err)

		assert.Equal(t, "answer\n", answer)
	})
}

func TestWithMockStdIn_ReturnsMultipleLines(t *testing.T) {
	utils.WithMockStdin(t, "answer1\nanswer2\n", func() {
		reader := bufio.NewReader(os.Stdin)
		answer1, err := reader.ReadString('\n')
		require.NoError(t, err)
		assert.Equal(t, "answer1\n", answer1)
		answer2, err := reader.ReadString('\n')
		require.NoError(t, err)
		assert.Equal(t, "answer2\n", answer2)
	})
}

func TestWithMockStdIn_AddsNewlineIfMissing(t *testing.T) {
	utils.WithMockStdin(t, "answer1", func() {
		reader := bufio.NewReader(os.Stdin)
		answer1, err := reader.ReadString('\n')
		require.NoError(t, err)
		assert.Equal(t, "answer1\n", answer1)
	})
}

func TestWithMockStdio_ReturnsInputSentToStdout(t *testing.T) {
	res := utils.WithMockStdio(t, "World", func() {
		reader := bufio.NewReader(os.Stdin)
		answer, err := reader.ReadString('\n')
		require.NoError(t, err)
		fmt.Fprintf(os.Stdout, "Hello %s", answer)
	})
	assert.Equal(t, "Hello World\n", res)
}
