package config

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestAzureNamifyReturnsHyphenatedStr(t *testing.T) {
	if err := LoadConfig(pflag.NewFlagSet("config", 1), false, "./testdata/config_test", "", true); err != nil {
		assert.Fail(t, "LoadConfig failed")
	}
	assert.Equal(t, namify("key", 20), "test-key-kekprt")
}

func TestAzureNamifyTruncatesNamesThatAreTooLong(t *testing.T) {
	if err := LoadConfig(pflag.NewFlagSet("config", 1), false, "./testdata/config_test", "", true); err != nil {
		assert.Fail(t, "LoadConfig failed")
	}
	namified := namify("this-is-a-really-long-key-that-should-be-truncated", 20)
	expected := "test-this-is-a-reall"
	assert.Equal(t, namified, expected)
}
