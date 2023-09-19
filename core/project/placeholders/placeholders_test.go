package placeholders

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClearString_ReturnsValidNames(t *testing.T) {
	type caseTest struct {
		input    string
		expected string
	}
	for _, test := range []caseTest{
		{"valid", "valid"},
		{"validCaps", "validCaps"},
		{"valid-Dash", "valid-Dash"},
		{"valid_Underscore", "valid_Underscore"},
		{"invalid&", "invalid"},
		{"invalid12", "invalid"},
		{"invalid£", "invalid"},
		{"mixed-Name_with&%%12", "mixed-Name_with"},
	} {
		actual := clearName(test.input)
		assert.Equal(t, test.expected, actual)
	}
}
