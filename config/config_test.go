package config

import (
	"bufio"
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vision-cli/vision/utils"
)

func TestLoadConfig_SetsViperDictionary(t *testing.T) {
	if err := LoadConfig(pflag.NewFlagSet("config", 1), false, "./testdata/config_test", ""); err != nil {
		assert.Fail(t, "LoadConfig failed")
	}
	assert.Equal(t, ProjectName(), "test")
	assert.Equal(t, Remote(), "remote")
	assert.Equal(t, TemplateVersion(), "v1")
	assert.Equal(t, Deployment(), "standalone-graphql")
}

func TestLoadConfig_NoConfigNotSilentUserDoesntCreate_ReturnError(t *testing.T) {
	utils.WithMockStdio(t, "n\n", func() {
		old := stat
		defer func() { stat = old }()
		stat = func(name string) (os.FileInfo, error) {
			return nil, os.ErrNotExist
		}
		if err := LoadConfig(pflag.NewFlagSet("config", 1), false, "./testdata/config_test", ""); err == nil {
			assert.Fail(t, "Test should fail if user doesn't want to create config")
		}
	})
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

func TestLoadDefaultConfig_SilentWithAllDefaults_SetsAllConfig(t *testing.T) {
	old := v
	defer func() { v = old }()
	v = Persist(NewMockPersist())
	flagSet := pflag.NewFlagSet("config", 1)
	flagSet.String(FlagRemote, "github.com/mycompany", "")
	err := loadDefaultConfig(flagSet, true, "configfile", "projectname", bufio.NewReader(os.Stdin))
	require.NoError(t, err)
	assert.Equal(t, Remote(), "github.com/mycompany")
	assert.Equal(t, ProjectName(), "projectname")
	assert.Equal(t, defaultBranch, Branch())
	assert.Equal(t, defaultTempalateVersion, TemplateVersion())
}

func TestLoadDefaultConfig_SilentWithMissingDefaults_ReturnsError(t *testing.T) {
	old := v
	defer func() { v = old }()
	v = Persist(NewMockPersist())
	flagSet := pflag.NewFlagSet("config", 1)
	err := loadDefaultConfig(flagSet, true, "configfile", "projectname", bufio.NewReader(os.Stdin))
	require.Error(t, err)
}

func TestLoadDefaultConfig_NotSilent_SetsAllConfig(t *testing.T) {
	old := v
	defer func() { v = old }()
	m := NewMockPersist()
	v = Persist(m)

	utils.WithMockStdio(t, "jarvis\ngithub.com/starkindustries\n\n\n\n\n\n\n\n\n\n\n", func() {
		reader := bufio.NewReader(os.Stdin)
		flagSet := pflag.NewFlagSet("config", 1)
		err := loadDefaultConfig(flagSet, false, "configfile", "projectname", reader)
		require.NoError(t, err)
		assert.Equal(t, Remote(), "github.com/starkindustries")
		assert.Equal(t, ProjectName(), "jarvis")
	})

}
