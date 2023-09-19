package tmpl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type tmplTest struct {
	inputText   string
	inputTokens interface{}
	expected    string
}

func TestTmplToString(t *testing.T) {
	type Tokens struct {
		Token1 string
	}
	for _, test := range []tmplTest{
		{"test{{.Token1}}", Tokens{"ThisToken"}, "testThisToken"},
		{"test{{.Token1 | Snake}}", Tokens{"ThisToken"}, "testthis_token"},
	} {
		actual, err := TmplToString(test.inputText, test.inputTokens)
		assert.Nil(t, err)
		assert.Equal(t, test.expected, actual)
	}
}
