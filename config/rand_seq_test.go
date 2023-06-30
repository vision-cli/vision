package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandSeq_CreatesCorrectLength(t *testing.T) {
	s4 := randSeq(4)
	assert.Equal(t, len(s4), 4)
	s7 := randSeq(7)
	assert.Equal(t, len(s7), 7)
}
