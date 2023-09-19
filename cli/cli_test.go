package cli_test

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/vision-cli/vision/common/execute"

	"github.com/stretchr/testify/assert"
	"github.com/vision-cli/vision/cli"
	"github.com/vision-cli/vision/common/mocks"
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
	output := mocks.WithMockStdout(t, func() {
		cli.Confirmed(reader, message)
	})
	assert.Equal(t, message+" (Y/n): ", output)
}

func TestInput_WhenAnswered_ReturnsAnswer(t *testing.T) {
	mocks.WithMockStdio(t, "answer\n", func() {
		reader := bufio.NewReader(os.Stdin)
		assert.Equal(t, "answer", cli.Input(reader, "Question", "", true))
	})
}

func TestInput_ReturnsAnswerWhenDefaultAndAnswerPresent(t *testing.T) {
	mocks.WithMockStdio(t, "answer\n", func() {
		reader := bufio.NewReader(os.Stdin)
		assert.Equal(t, "answer", cli.Input(reader, "Question", "default", true))
	})
}

func TestInput_ReturnsDefaultWhenPresentAndNoAnswer(t *testing.T) {
	mocks.WithMockStdio(t, "\n", func() {
		reader := bufio.NewReader(os.Stdin)
		assert.Equal(t, "default", cli.Input(reader, "Question", "default", true))
	})
}

func TestInput_DontHaveToAnswerNonMandatoryWithoutDefault(t *testing.T) {
	mocks.WithMockStdio(t, "\n", func() {
		reader := bufio.NewReader(os.Stdin)
		assert.Equal(t, "", cli.Input(reader, "Question", "", false))
	})
}

func TestInput_MustAnswerMandatory(t *testing.T) {
	mocks.WithMockStdio(t, "\nanswer\n", func() {
		reader := bufio.NewReader(os.Stdin)
		assert.Equal(t, "answer", cli.Input(reader, "Question", "", true))
	})
}

func TestInputWithValidation_WithValidResult_ReturnsResult(t *testing.T) {
	mocks.WithMockStdio(t, "answer\n", func() {
		reader := bufio.NewReader(os.Stdin)
		te := mocks.NewMockExecutor()
		validationFunc := func(s string, e execute.Executor) (bool, string) {
			return true, s
		}
		assert.Equal(t, "answer", cli.InputWithValidation(reader, "Question", "", true, validationFunc, &te))
	})
}

func TestInputWithValidation_WithInValidResult_AsksConfirmThenReturnsResult(t *testing.T) {
	mocks.WithMockStdio(t, "answer\ny\n", func() {
		reader := bufio.NewReader(os.Stdin)
		te := mocks.NewMockExecutor()
		validationFunc := func(s string, e execute.Executor) (bool, string) {
			return false, s
		}
		assert.Equal(t, "answer", cli.InputWithValidation(reader, "Question", "", true, validationFunc, &te))
	})
}

func testIsConfirmed_ForEachInput_Returns(t *testing.T, inputs []string, expected bool) {
	t.Helper()
	for _, input := range inputs {
		mocks.WithMockStdio(t, fmt.Sprintf("%s\n", input), func() {
			reader := bufio.NewReader(os.Stdin)
			assert.Equal(t, expected, cli.Confirmed(reader, "Do you confirm?"))
		})
	}
}
