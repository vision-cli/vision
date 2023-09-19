package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/vision-cli/common/execute"
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
			log.Error(err.Error())
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
		log.Warn("invalid input: %s", msg)
		if Confirmed(reader, "Do you still want to continue?") {
			return s
		}
	}
}
