package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/vision-cli/vision/execute"
)

const (
	fatal   = "FATAL!  "
	warning = "WARNING!"
	hint    = "HINT    "
	info    = "INFO    "
)

type Colour int

const (
	Black Colour = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	Reset Colour = 0
)

// Confirmed returns input from the user in answer to the yes/no question.
func Confirmed(reader *bufio.Reader, question string) bool {
	fmt.Printf("%s (Y/n): ", question)

	answer, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	s := strings.TrimSpace(strings.ToLower(answer))
	return s == "y" || s == "yes" || s == ""
}

// Input returns input from the user in answer to the question.
func Input(reader *bufio.Reader, question, dflt string, mandatory bool) string {
	for {
		if dflt == "" {
			fmt.Printf("%s : ", question)
		} else {
			fmt.Printf("%s (default %s): ", question, dflt)
		}

		var s string
		answer, err := reader.ReadString('\n')
		if err != nil {
			Fatalf(err.Error())
		}
		s = strings.TrimSpace(answer)
		// Combinations
		// mandatory  answered  default  result
		// true       true      true     answer
		// true       true      false    answer
		// true       false     true     dflt
		// true       false     false    retry
		// false      true      true     answer
		// false      true      false    answer
		// false      false     true     dflt
		// false      false     false    blank
		if s != "" {
			return s
		}
		if dflt != "" {
			return dflt
		}
		if !mandatory {
			return ""
		}
		//retry
		fmt.Println("please provide an answer")
	}
}

// Input returns input from the user in answer to the question.
func InputWithValidation(
	reader *bufio.Reader,
	question,
	dflt string,
	mandatory bool,
	validation func(string, execute.Executor) (bool, string),
	executor execute.Executor) string {
	for {
		s := Input(reader, question, dflt, mandatory)
		valid, msg := validation(s, executor)
		if valid {
			return s
		}
		Warningf("invalid input: %s", msg)
		if Confirmed(reader, "Do you still want to continue?") {
			return s
		}
	}
}

// Warningf prints formatted message with red "WARNING!" prefix to stderr.
func Warningf(message string, a ...any) {
	fmt.Fprintf(os.Stderr, "%s %s\n", Highlight(warning, Yellow), fmt.Sprintf(message, a...))
}

// Warningf prints formatted message with red "WARNING!" prefix to stderr.
func Fatalf(message string, a ...any) {
	fmt.Fprintf(os.Stderr, "%s %s\n", Highlight(fatal, Red), fmt.Sprintf(message, a...))
	os.Exit(1)
}

// Hintf prints formatted message with yellow "HINT" prefix to stdout.
func Hintf(message string, a ...any) {
	fmt.Fprintf(os.Stdout, "%s %s\n", Highlight(hint, Yellow), fmt.Sprintf(message, a...))
}

// Infof prints formatted message with gree "INFO" prefix to stdout.
func Infof(message string, a ...any) {
	fmt.Fprintf(os.Stdout, "%s %s\n", Highlight(info, Green), fmt.Sprintf(message, a...))
}

func Highlight(s string, colour Colour) string {
	return esc(colour) + s + esc(Reset)
}

func esc(colour Colour) string {
	return fmt.Sprintf("\033[%dm", colour)
}
