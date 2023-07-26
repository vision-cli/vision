package config

import (
	"bufio"
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vision-cli/common/mocks"
)

const (
	stdinPass = "jarvis\ngithub.com/starkindustries\n\n\ngrc.io/stark\n\n\n\n\n\n\n\n"
)

func TestLoadConfig_SetsViperDictionary(t *testing.T) {
	if err := LoadConfig(pflag.NewFlagSet("config", 1), false, "./testdata/config_test", "", true); err != nil {
		assert.Fail(t, "LoadConfig failed")
	}
	assert.Equal(t, ProjectName(), "test")
	assert.Equal(t, Remote(), "remote")
	assert.Equal(t, TemplateVersion(), "v1")
	assert.Equal(t, Deployment(), "standalone-graphql")
	assert.Equal(t, v.GetString(uniqueStr), "kekprt")
}

func TestLoadConfig_NoConfigNotSilentUserDoesntCreate_ReturnError(t *testing.T) {
	mocks.WithMockStdio(t, "n\n", func() {
		old := stat
		defer func() { stat = old }()
		stat = func(name string) (os.FileInfo, error) {
			return nil, os.ErrNotExist
		}
		if err := LoadConfig(pflag.NewFlagSet("config", 1), false, "./testdata/config_test", "", true); err == nil {
			assert.Fail(t, "Test should fail if user doesn't want to create config")
		}
	})
}

func TestLoadConfig_NoConfigSilent_FailsDueToMissingRemote(t *testing.T) {
	oldv := v
	defer func() { v = oldv }()
	v = Persist(NewMockPersist())

	oldstat := stat
	defer func() { stat = oldstat }()
	stat = func(name string) (os.FileInfo, error) {
		return nil, os.ErrNotExist
	}

	if err := LoadConfig(pflag.NewFlagSet("config", 1), true, "./testdata/config_test", "", true); err == nil {
		assert.Fail(t, "Test should fail because user has not provided a remote")
	}
}

func TestLoadConfig_NoConfigSilentWithRemote_CreatesDefaultConfig(t *testing.T) {
	oldv := v
	defer func() { v = oldv }()
	v = Persist(NewMockPersist())

	oldstat := stat
	defer func() { stat = oldstat }()
	stat = func(name string) (os.FileInfo, error) {
		return nil, os.ErrNotExist
	}

	flagSet := pflag.NewFlagSet("config", 1)
	flagSet.String(FlagRemote, "github.com/mycompany", "")
	flagSet.String(FlagRegistry, "gcr.io/mycompany", "")

	err := LoadConfig(flagSet, true, "./testdata/config_test", "", true)
	require.NoError(t, err)
	assert.Equal(t, Remote(), "github.com/mycompany")
	assert.Equal(t, TemplateVersion(), "v1")
	assert.Equal(t, Deployment(), "standalone-graphql")
	assert.Equal(t, len(v.GetString(uniqueStr)), 6)
}

func TestLoadConfig_NoConfigNotSilent_CreatesDefaultConfig(t *testing.T) {
	oldv := v
	defer func() { v = oldv }()
	v = Persist(NewMockPersist())

	oldstat := stat
	defer func() { stat = oldstat }()
	stat = func(name string) (os.FileInfo, error) {
		return nil, os.ErrNotExist
	}

	mocks.WithMockStdio(t, "y\n"+stdinPass, func() {
		err := LoadConfig(pflag.NewFlagSet("config", 1), false, "./testdata/config_test", "", true)
		require.NoError(t, err)
		assert.Equal(t, Remote(), "github.com/starkindustries")
		assert.Equal(t, TemplateVersion(), "v1")
		assert.Equal(t, Deployment(), "standalone-graphql")
		assert.Equal(t, len(v.GetString(uniqueStr)), 6)
	})
}

func TestLoadConfig_NoConfigNotSilentWithFlag_DoesntPromptAndSetsFlagValue(t *testing.T) {
	oldv := v
	defer func() { v = oldv }()
	v = Persist(NewMockPersist())

	oldstat := stat
	defer func() { stat = oldstat }()
	stat = func(name string) (os.FileInfo, error) {
		return nil, os.ErrNotExist
	}

	flagSet := pflag.NewFlagSet("config", 1)
	flagSet.String(FlagRemote, "github.com/ironman", "")

	mocks.WithMockStdio(t, "y\njarvis\n\n\ngrc.io/stark\n\n\n\n\n\n\n\n", func() {
		err := LoadConfig(flagSet, false, "./testdata/config_test", "", true)
		require.NoError(t, err)
		assert.Equal(t, Remote(), "github.com/ironman")
		assert.Equal(t, TemplateVersion(), "v1")
		assert.Equal(t, Deployment(), "standalone-graphql")
		assert.Equal(t, len(v.GetString(uniqueStr)), 6)
	})
}

func TestGenericSetter_SetsDefaultWhenSilent(t *testing.T) {
	genericSetter(bufio.NewReader(os.Stdin), "", "test", true, func(val string) {
		assert.Equal(t, val, "test")
	})
}

func TestGenericSetter_PromptsAndSetsUserInputWhenNotSilent(t *testing.T) {
	mocks.WithMockStdio(t, "anothertest\n", func() {
		reader := bufio.NewReader(os.Stdin)
		genericSetter(reader, "", "test", false, func(val string) {
			assert.Equal(t, val, "anothertest")
		})
	})
}

func TestMustSetWithFlag_PrioritisesFlag(t *testing.T) {
	flagSet := pflag.NewFlagSet("config", 1)
	flagSet.String("exampleflag", "flagoverride", "")
	err := mustSetWithFlag(bufio.NewReader(os.Stdin), "", "defaultval", false, flagSet, "exampleflag", func(val string) {
		assert.Equal(t, val, "flagoverride")
	})
	require.NoError(t, err)
}

func TestMustSetWithFlag_UsesDefaultIfSilentAndNoFlag(t *testing.T) {
	flagSet := pflag.NewFlagSet("config", 1)
	flagSet.String("exampleflag", "", "")
	err := mustSetWithFlag(bufio.NewReader(os.Stdin), "", "defaultval", true, flagSet, "exampleflag", func(val string) {
		assert.Equal(t, val, "defaultval")
	})
	require.NoError(t, err)
}

func TestMustSetWithFlag_PromptsWithDefaultIfNotSilentAndNoFlag(t *testing.T) {
	mocks.WithMockStdio(t, "anothertest\n", func() {
		reader := bufio.NewReader(os.Stdin)
		flagSet := pflag.NewFlagSet("config", 1)
		flagSet.String("exampleflag", "", "")
		err := mustSetWithFlag(reader, "", "defaultval", false, flagSet, "exampleflag", func(val string) {
			assert.Equal(t, val, "anothertest")
		})
		require.NoError(t, err)
	})
}

func TestLoadDefaultConfig_SilentWithAllDefaults_SetsAllConfig(t *testing.T) {
	old := v
	defer func() { v = old }()
	v = Persist(NewMockPersist())
	flagSet := pflag.NewFlagSet("config", 1)
	flagSet.String(FlagRemote, "github.com/mycompany", "")
	flagSet.String(FlagRegistry, "gcr.io/mycompany", "")
	err := loadDefaultConfig(flagSet, true, "configfile", "projectname", bufio.NewReader(os.Stdin))
	require.NoError(t, err)
	assert.Equal(t, Remote(), "github.com/mycompany")
	assert.Equal(t, ProjectName(), "projectname")
	assert.Equal(t, defaultBranch, Branch())
	assert.Equal(t, defaultTemplateVersion, TemplateVersion())
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
	v = Persist(NewMockPersist())

	mocks.WithMockStdio(t, stdinPass, func() {
		reader := bufio.NewReader(os.Stdin)
		flagSet := pflag.NewFlagSet("config", 1)
		err := loadDefaultConfig(flagSet, false, "configfile", "projectname", reader)
		require.NoError(t, err)
		assert.Equal(t, Remote(), "github.com/starkindustries")
		assert.Equal(t, ProjectName(), "jarvis")
	})

}
