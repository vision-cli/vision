package placeholders

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type caseTest struct {
	input    string
	expected string
}

func TestPascalCase(t *testing.T) {
	for _, test := range []caseTest{
		{"word", "Word"},
		{"two-words", "TwoWords"},
		{"TwoWords", "TwoWords"},
		{"Only-% alpha \nNumeric!123", "OnlyAlphaNumeric123"},
		{"123abC", "123AbC"},
		{"can'tBe-apostrophe", "CantBeApostrophe"},
	} {
		actual := Pascal(test.input)
		assert.Equal(t, test.expected, actual)
	}
}

func TestCamelCase(t *testing.T) {
	for _, test := range []caseTest{
		{"Word", "word"},
		{"two-words", "twoWords"},
		{"twoWords", "twoWords"},
		{"Only-% alpha \nNumeric!123", "onlyAlphaNumeric123"},
		{"123abC", "123AbC"},
		{"can'tBe-apostrophe", "cantBeApostrophe"},
	} {
		actual := Camel(test.input)
		assert.Equal(t, test.expected, actual)
	}
}

func TestKebabCase(t *testing.T) {
	for _, test := range []caseTest{
		{"word", "word"},
		{"TwoWords", "two-words"},
		{"two_words", "two-words"},
		{"Only-% alpha \nNumeric!123", "only-alpha-numeric-123"},
		{"123abC", "123-ab-c"},
		{"can'tBe-apostrophe", "cant-be-apostrophe"},
	} {
		actual := Kebab(test.input)
		assert.Equal(t, test.expected, actual)
	}
}

func TestSnakeCase(t *testing.T) {
	for _, test := range []caseTest{
		{"word", "word"},
		{"TwoWords", "two_words"},
		{"two_words", "two_words"},
		{"Only-% alpha \nNumeric!123", "only_alpha_numeric_123"},
		{"123abC", "123_ab_c"},
		{"can'tBe-apostrophe", "cant_be_apostrophe"},
	} {
		actual := Snake(test.input)
		assert.Equal(t, test.expected, actual)
	}
}
