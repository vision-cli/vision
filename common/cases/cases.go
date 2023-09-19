package cases

import (
	"regexp"
	"strings"
)

var distinctWordsExp = regexp.MustCompile(`[a-zA-Z][a-z]*|[\d]+`)

func Pascal(s string) string {
	words := toWords(s)
	for i, word := range words {
		words[i] = strings.ToUpper(word[:1]) + word[1:]
	}
	return strings.Join(words, "")
}

func Camel(s string) string {
	words := toLowerWords(s)
	for i, word := range words {
		if i != 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, "")
}

func Snake(s string) string {
	words := toLowerWords(s)
	return strings.Join(words, "_")
}

func Kebab(s string) string {
	words := toLowerWords(s)
	return strings.Join(words, "-")
}

func toWords(s string) []string {
	ss := strings.ReplaceAll(s, "'", "")
	return distinctWordsExp.FindAllString(ss, -1)
}

func toLowerWords(s string) []string {
	words := toWords(s)
	for i, word := range words {
		words[i] = strings.ToLower(word)
	}
	return words
}
