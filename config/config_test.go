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
	assert.Equal(t, "test", ProjectName())
	assert.Equal(t, "remote", Remote())
	assert.Equal(t, "v1", TemplateVersion())
	assert.Equal(t, "standalone-graphql", Deployment())
	assert.Equal(t, "kekprt", v.GetString(uniqueStr))
	assert.Equal(t, "master", Branch())
}

func TestLoadConfig_NoConfigRequiredConfigFileExists_LoadsConfigFile(t *testing.T) {
	if err := LoadConfig(pflag.NewFlagSet("config", 1), false, "./testdata/config_test", "", false); err != nil {
		assert.Fail(t, "LoadConfig failed")
	}
	assert.Equal(t, "test", ProjectName())
	assert.Equal(t, "remote", Remote())
	assert.Equal(t, "kekprt", v.GetString(uniqueStr))
	assert.Equal(t, "master", Branch())
}

func TestLoadConfig_NoConfigRequiredConfigFileDoesntExist_LoadsDefaultValues(t *testing.T) {
	oldv := v
	defer func() { v = oldv }()
	v = Persist(NewMockPersist())

	oldstat := stat
	defer func() { stat = oldstat }()
	stat = func(name string) (os.FileInfo, error) {
		return nil, os.ErrNotExist
	}

	if err := LoadConfig(pflag.NewFlagSet("config", 1), false, "./testdata/config_test", "", false); err != nil {
		assert.Fail(t, "LoadConfig failed")
	}
	assert.Equal(t, "", ProjectName())
	assert.Equal(t, "", Remote())
	assert.Equal(t, "master", Branch())
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
	assert.Equal(t, "github.com/mycompany", Remote())
	assert.Equal(t, "v1", TemplateVersion())
	assert.Equal(t, "standalone-graphql", Deployment())
	assert.Equal(t, 6, len(v.GetString(uniqueStr)))
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
		assert.Equal(t, "github.com/starkindustries", Remote())
		assert.Equal(t, "v1", TemplateVersion())
		assert.Equal(t, "standalone-graphql", Deployment())
		assert.Equal(t, 6, len(v.GetString(uniqueStr)))
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
		assert.Equal(t, "github.com/ironman", Remote())
		assert.Equal(t, "v1", TemplateVersion())
		assert.Equal(t, "standalone-graphql", Deployment())
		assert.Equal(t, 6, len(v.GetString(uniqueStr)))
	})
}

func TestGenericSetter_SetsDefaultWhenSilent(t *testing.T) {
	genericSetter(bufio.NewReader(os.Stdin), "", "test", true, func(val string) {
		assert.Equal(t, "test", val)
	})
}

func TestGenericSetter_PromptsAndSetsUserInputWhenNotSilent(t *testing.T) {
	mocks.WithMockStdio(t, "anothertest\n", func() {
		reader := bufio.NewReader(os.Stdin)
		genericSetter(reader, "", "test", false, func(val string) {
			assert.Equal(t, "anothertest", val)
		})
	})
}

func TestMustSetWithFlag_PrioritisesFlag(t *testing.T) {
	flagSet := pflag.NewFlagSet("config", 1)
	flagSet.String("exampleflag", "flagoverride", "")
	err := mustSetWithFlag(bufio.NewReader(os.Stdin), "", "defaultval", false, flagSet, "exampleflag", func(val string) {
		assert.Equal(t, "flagoverride", val)
	})
	require.NoError(t, err)
}

func TestMustSetWithFlag_UsesDefaultIfSilentAndNoFlag(t *testing.T) {
	flagSet := pflag.NewFlagSet("config", 1)
	flagSet.String("exampleflag", "", "")
	err := mustSetWithFlag(bufio.NewReader(os.Stdin), "", "defaultval", true, flagSet, "exampleflag", func(val string) {
		assert.Equal(t, "defaultval", val)
	})
	require.NoError(t, err)
}

func TestMustSetWithFlag_PromptsWithDefaultIfNotSilentAndNoFlag(t *testing.T) {
	mocks.WithMockStdio(t, "anothertest\n", func() {
		reader := bufio.NewReader(os.Stdin)
		flagSet := pflag.NewFlagSet("config", 1)
		flagSet.String("exampleflag", "", "")
		err := mustSetWithFlag(reader, "", "defaultval", false, flagSet, "exampleflag", func(val string) {
			assert.Equal(t, "anothertest", val)
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
	assert.Equal(t, "github.com/mycompany", Remote())
	assert.Equal(t, "projectname", ProjectName())
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
		assert.Equal(t, "github.com/starkindustries", Remote())
		assert.Equal(t, "jarvis", ProjectName())
	})

}
