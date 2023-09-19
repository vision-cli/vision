package file

import (
	"bufio"
	"os"
	"strings"
)

var ToLines = toLines
var FromLines = fromLines

// ToLines return the contents of the file at the specified path as a slice of strings.
func toLines(path string) ([]string, error) {
	lines := []string{}

	f, err := Osopen(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// Inserts strings into the slice immediately after the line containing pattern.
// The resulting slice is returned. If no lines contain pattern, the original slice is returned.
func InsertIntoLines(lines []string, pattern string, insert ...string) []string {
	for i, line := range lines {
		if strings.Contains(line, pattern) {
			return append(lines[:i+1], append(insert, lines[i+1:]...)...)
		}
	}
	return lines
}

// ensures all have a newline character suffix.
func clean(lines []string) []string {
	for i, line := range lines {
		if !strings.HasSuffix(line, "\n") {
			lines[i] = line + "\n"
		}
	}
	return lines
}

// FromLines writes lines to the file at the specified path, creating the file if none exists.
// Existing files are truncated before writing.
func fromLines(path string, lines []string) error {
	fileContents := strings.Join(clean(lines), "")
	return os.WriteFile(path, []byte(fileContents), os.ModePerm)
}
