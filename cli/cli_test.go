package cli_test

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vision-cli/vision/cli"
	"github.com/vision-cli/vision/execute"
)

func TestIsConfirmed_WhenAnsweredYesOrBlank_ReturnsTrue(t *testing.T) {
	confirmingInputs := []string{"y", "yes", "yES", "Y", ""}
	testIsConfirmed_ForEachInput_Returns(t, confirmingInputs, true)
}

func TestIsConfirmed_IgnoresWhitespace(t *testing.T) {
	confirmingInputs := []string{" y", "\tyes    "}
	testIsConfirmed_ForEachInput_Returns(t, confirmingInputs, true)
}

func TestIsConfirmed_WhenAnsweredNotYes_ReturnsFalse(t *testing.T) {
	rejectingInputs := []string{"no", " N", "idk"}
	testIsConfirmed_ForEachInput_Returns(t, rejectingInputs, false)
}

func TestIsConfirmed_PrintsMessageToStdout(t *testing.T) {
	message := "Do you confirm?"
	reader := bufio.NewReader(os.Stdin)
	output := withMockStdout(t, func() {
		cli.Confirmed(reader, message)
	})
	assert.Equal(t, message+" (Y/n): ", output)
}

func TestInput_WhenAnswered_ReturnsAnswer(t *testing.T) {
	withMockStdio(t, "answer\n", func() {
		reader := bufio.NewReader(os.Stdin)
		assert.Equal(t, "answer", cli.Input(reader, "Question", "", true))
	})
}

func TestInput_ReturnsAnswerWhenDefaultAndAnswerPresent(t *testing.T) {
	withMockStdio(t, "answer\n", func() {
		reader := bufio.NewReader(os.Stdin)
		assert.Equal(t, "answer", cli.Input(reader, "Question", "default", true))
	})
}

func TestInput_ReturnsDefaultWhenPresentAndNoAnswer(t *testing.T) {
	withMockStdio(t, "\n", func() {
		reader := bufio.NewReader(os.Stdin)
		assert.Equal(t, "default", cli.Input(reader, "Question", "default", true))
	})
}

func TestInput_DontHaveToAnswerNonMandatoryWithoutDefault(t *testing.T) {
	withMockStdio(t, "\n", func() {
		reader := bufio.NewReader(os.Stdin)
		assert.Equal(t, "", cli.Input(reader, "Question", "", false))
	})
}

func TestInput_MustAnswerMandatory(t *testing.T) {
	withMockStdio(t, "\nanswer\n", func() {
		reader := bufio.NewReader(os.Stdin)
		assert.Equal(t, "answer", cli.Input(reader, "Question", "", true))
	})
}

func TestInputWithValidation_WithValidResult_ReturnsResult(t *testing.T) {
	withMockStdio(t, "answer\n", func() {
		reader := bufio.NewReader(os.Stdin)
		te := execute.NewMockExecutor()
		validationFunc := func(s string, e execute.Executor) (bool, string) {
			return true, s
		}
		assert.Equal(t, "answer", cli.InputWithValidation(reader, "Question", "", true, validationFunc, te))
	})
}

func TestInputWithValidation_WithInValidResult_AsksConfirmThenReturnsResult(t *testing.T) {
	withMockStdio(t, "answer\ny\n", func() {
		reader := bufio.NewReader(os.Stdin)
		te := execute.NewMockExecutor()
		validationFunc := func(s string, e execute.Executor) (bool, string) {
			return false, s
		}
		assert.Equal(t, "answer", cli.InputWithValidation(reader, "Question", "", true, validationFunc, te))
	})
}

func TestWithMockStdIn_ReturnsSingleLine(t *testing.T) {
	withMockStdin(t, "answer\n", func() {
		reader := bufio.NewReader(os.Stdin)
		answer, err := reader.ReadString('\n')
		require.NoError(t, err)

		assert.Equal(t, "answer\n", answer)
	})
}

func TestWithMockStdIn_ReturnsMultipleLines(t *testing.T) {
	withMockStdin(t, "answer1\nanswer2\n", func() {
		reader := bufio.NewReader(os.Stdin)
		answer1, err := reader.ReadString('\n')
		require.NoError(t, err)
		assert.Equal(t, "answer1\n", answer1)
		answer2, err := reader.ReadString('\n')
		require.NoError(t, err)
		assert.Equal(t, "answer2\n", answer2)
	})
}

func testIsConfirmed_ForEachInput_Returns(t *testing.T, inputs []string, expected bool) {
	t.Helper()
	for _, input := range inputs {
		withMockStdio(t, fmt.Sprintf("%s\n", input), func() {
			reader := bufio.NewReader(os.Stdin)
			assert.Equal(t, expected, cli.Confirmed(reader, "Do you confirm?"))
		})
	}
}

func withMockStdio(t *testing.T, input string, test func()) string {
	t.Helper()
	return withMockStdout(t, func() {
		withMockStdin(t, input, test)
	})
}

func withMockStdout(t *testing.T, test func()) string {
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

func withMockStdin(t *testing.T, input string, test func()) {
	t.Helper()
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
