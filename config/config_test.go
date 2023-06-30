package config

import (
	"bufio"
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/vision-cli/vision/utils"
)

func TestRandSeq_CreatesCorrectLength(t *testing.T) {
	s4 := randSeq(4)
	assert.Equal(t, len(s4), 4)
	s7 := randSeq(7)
	assert.Equal(t, len(s7), 7)
}

func TestLoadConfig_SetsViperDictionary(t *testing.T) {
	if err := LoadConfig(pflag.NewFlagSet("config", 1), false, "./testdata/config_test", ""); err != nil {
		assert.Fail(t, "LoadConfig failed")
	}
	assert.Equal(t, ProjectName(), "test")
	assert.Equal(t, Remote(), "remote")
	assert.Equal(t, TemplateVersion(), "v1")
	assert.Equal(t, Deployment(), "standalone-graphql")
}

func TestGenericSetter_SetsDefaultWhenSilent(t *testing.T) {
	genericSetter(bufio.NewReader(os.Stdin), "", "test", true, func(val string) {
		assert.Equal(t, val, "test")
	})
}

func TestGenericSetter_PromptsAndSetsUserInputWhenNotSilent(t *testing.T) {
	utils.WithMockStdio(t, "anothertest\n", func() {
		reader := bufio.NewReader(os.Stdin)
		genericSetter(reader, "", "test", false, func(val string) {
			assert.Equal(t, val, "anothertest")
		})
	})
}

func TestMustSetWithFlag_PrioritisesFlag(t *testing.T) {
	flagSet := pflag.NewFlagSet("config", 1)
	flagSet.String("exampleflag", "flagoverride", "")
	mustSetWithFlag(bufio.NewReader(os.Stdin), "", "defaultval", false, flagSet, "exampleflag", func(val string) {
		assert.Equal(t, val, "flagoverride")
	})
}

func TestMustSetWithFlag_UsesDefaultIfSilentAndNoFlag(t *testing.T) {
	flagSet := pflag.NewFlagSet("config", 1)
	flagSet.String("exampleflag", "", "")
	mustSetWithFlag(bufio.NewReader(os.Stdin), "", "defaultval", true, flagSet, "exampleflag", func(val string) {
		assert.Equal(t, val, "defaultval")
	})
}

func TestMustSetWithFlag_PromptsWithDefaultIfNotSilentAndNoFlag(t *testing.T) {
	utils.WithMockStdio(t, "anothertest\n", func() {
		reader := bufio.NewReader(os.Stdin)
		flagSet := pflag.NewFlagSet("config", 1)
		flagSet.String("exampleflag", "", "")
		mustSetWithFlag(reader, "", "defaultval", false, flagSet, "exampleflag", func(val string) {
			assert.Equal(t, val, "anothertest")
		})
	})
}
